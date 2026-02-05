package image

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createTestImageData() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

type mockImageService struct {
	images map[uint]*Image
	nextID uint
}

func newMockImageService() *mockImageService {
	return &mockImageService{
		images: make(map[uint]*Image),
		nextID: 1,
	}
}

func (m *mockImageService) SaveImage(image *Image) (*Image, error) {
	image.ID = m.nextID
	m.nextID++
	m.images[image.ID] = image
	return image, nil
}

func (m *mockImageService) GetImageByID(id uint) (*Image, error) {
	if img, ok := m.images[id]; ok {
		return img, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockImageService) ListImagesByUser(userID uint) ([]Image, error) {
	var images []Image
	for _, img := range m.images {
		if img.UserID == userID {
			images = append(images, *img)
		}
	}
	return images, nil
}

func (m *mockImageService) DeleteImage(id uint) error {
	if _, ok := m.images[id]; ok {
		delete(m.images, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

type mockUploadService struct {
	uploadedFiles map[string]string
}

func newMockUploadService() *mockUploadService {
	return &mockUploadService{
		uploadedFiles: make(map[string]string),
	}
}

func (m *mockUploadService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	url := "http://mock.com/uploads/" + file.Filename
	m.uploadedFiles[file.Filename] = url
	return url, nil
}

func TestImageHandler_UploadImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockImgService := newMockImageService()
	mockUploadService := newMockUploadService()
	handler := NewImageHandler(mockImgService, mockUploadService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	fileWriter, _ := writer.CreateFormFile("file", "test.jpg")
	fileWriter.Write(createTestImageData())

	// Add alt_text
	writer.WriteField("alt_text", "Test image")

	writer.Close()

	c.Request = httptest.NewRequest("POST", "/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Set("user_id", uint(1))

	handler.UploadImage()(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Data not found in response")
	}

	url, ok := data["url"].(string)
	if !ok || url == "" {
		t.Error("URL should not be empty")
	}
}

func TestImageHandler_GetImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockImgService := newMockImageService()
	mockUploadService := newMockUploadService()
	handler := NewImageHandler(mockImgService, mockUploadService)

	// Pre-create image
	img := &Image{
		URL:     "http://example.com/image.jpg",
		AltText: "Test",
		UserID:  1,
	}
	mockImgService.SaveImage(img)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/images/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", uint(1))

	handler.GetImage()(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Data not found in response")
	}

	if data["id"].(float64) != 1 {
		t.Error("Wrong image ID")
	}
}

func TestImageHandler_GetImage_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockImgService := newMockImageService()
	mockUploadService := newMockUploadService()
	handler := NewImageHandler(mockImgService, mockUploadService)

	// Pre-create image for different user
	img := &Image{
		URL:     "http://example.com/image.jpg",
		AltText: "Test",
		UserID:  2,
	}
	mockImgService.SaveImage(img)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/images/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", uint(1)) // Different user

	handler.GetImage()(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestImageHandler_ListUserImages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockImgService := newMockImageService()
	mockUploadService := newMockUploadService()
	handler := NewImageHandler(mockImgService, mockUploadService)

	// Pre-create images
	img1 := &Image{URL: "url1", UserID: 1}
	img2 := &Image{URL: "url2", UserID: 1}
	mockImgService.SaveImage(img1)
	mockImgService.SaveImage(img2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/users/1/images", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}
	c.Set("user_id", uint(1))

	handler.ListUserImages()(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Data not found in response")
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 images, got %d", len(data))
	}
}

func TestImageHandler_DeleteImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockImgService := newMockImageService()
	mockUploadService := newMockUploadService()
	handler := NewImageHandler(mockImgService, mockUploadService)

	// Pre-create image
	img := &Image{URL: "url", UserID: 1}
	mockImgService.SaveImage(img)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("DELETE", "/images/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", uint(1))

	handler.DeleteImage()(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check deleted
	_, err := mockImgService.GetImageByID(1)
	if err == nil {
		t.Error("Image should be deleted")
	}
}
