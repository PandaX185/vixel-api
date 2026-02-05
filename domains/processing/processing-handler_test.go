package processing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockProcessingService struct {
	transformResults map[string]string
}

func newMockProcessingService() *mockProcessingService {
	return &mockProcessingService{
		transformResults: make(map[string]string),
	}
}

func (m *mockProcessingService) TransformImage(ctx context.Context, imageID string, dto TransformationDTO) (string, error) {
	// Mock transformation result
	url := "http://mock.com/transformed/" + imageID
	m.transformResults[imageID] = url
	return url, nil
}

func TestProcessingHandler_TransformImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := newMockProcessingService()
	handler := NewProcessingHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	dto := TransformationDTO{
		Resize: &ResizeDTO{
			Width:  100,
			Height: 100,
		},
	}
	jsonData, _ := json.Marshal(dto)

	c.Request = httptest.NewRequest("POST", "/images/123/transform", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	handler.TransformImage()(c)

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

	newURL, ok := data["new_image_url"].(string)
	if !ok || newURL == "" {
		t.Error("New image URL should not be empty")
	}

	if newURL != "http://mock.com/transformed/123" {
		t.Errorf("Expected URL 'http://mock.com/transformed/123', got '%s'", newURL)
	}
}

func TestProcessingHandler_TransformImage_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := newMockProcessingService()
	handler := NewProcessingHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/images/123/transform", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	handler.TransformImage()(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
