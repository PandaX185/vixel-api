package image

import (
	"testing"

	"vixel/domains/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupImageTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	db.AutoMigrate(&Image{}, &user.User{})
	return db
}

func TestImageService_SaveImage(t *testing.T) {
	db := setupImageTestDB(t)
	service := NewImageService(db)

	img := &Image{
		URL:     "http://example.com/image.jpg",
		AltText: "Test image",
		UserID:  1,
	}

	saved, err := service.SaveImage(img)
	if err != nil {
		t.Fatalf("SaveImage failed: %v", err)
	}

	if saved.ID == 0 {
		t.Error("Image ID should be set")
	}

	if saved.URL != "http://example.com/image.jpg" {
		t.Error("URL not saved correctly")
	}
}

func TestImageService_GetImageByID(t *testing.T) {
	db := setupImageTestDB(t)
	service := NewImageService(db)

	usr := &user.User{Username: "test", Email: "test@example.com", Password: "pass"}
	db.Create(usr)

	img := &Image{
		URL:     "http://example.com/image.jpg",
		AltText: "Test image",
		UserID:  usr.ID,
	}
	db.Create(img)

	retrieved, err := service.GetImageByID(img.ID)
	if err != nil {
		t.Fatalf("GetImageByID failed: %v", err)
	}

	if retrieved.ID != img.ID {
		t.Error("Retrieved wrong image")
	}

	if retrieved.User.ID != usr.ID {
		t.Error("User not preloaded correctly")
	}
}

func TestImageService_GetImageByID_NotFound(t *testing.T) {
	db := setupImageTestDB(t)
	service := NewImageService(db)

	_, err := service.GetImageByID(999)
	if err == nil {
		t.Error("Expected error for non-existent image")
	}
}

func TestImageService_ListImagesByUser(t *testing.T) {
	db := setupImageTestDB(t)
	service := NewImageService(db)

	usr := &user.User{Username: "test", Email: "test@example.com", Password: "pass"}
	db.Create(usr)

	img1 := &Image{URL: "url1", UserID: usr.ID}
	img2 := &Image{URL: "url2", UserID: usr.ID}
	db.Create(img1)
	db.Create(img2)

	usr2 := &user.User{Username: "test2", Email: "test2@example.com", Password: "pass"}
	db.Create(usr2)
	img3 := &Image{URL: "url3", UserID: usr2.ID}
	db.Create(img3)

	images, err := service.ListImagesByUser(usr.ID)
	if err != nil {
		t.Fatalf("ListImagesByUser failed: %v", err)
	}

	if len(images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(images))
	}
}

func TestImageService_DeleteImage(t *testing.T) {
	db := setupImageTestDB(t)
	service := NewImageService(db)

	img := &Image{URL: "url", UserID: 1}
	db.Create(img)

	err := service.DeleteImage(img.ID)
	if err != nil {
		t.Fatalf("DeleteImage failed: %v", err)
	}

	var count int64
	db.Model(&Image{}).Where("id = ?", img.ID).Count(&count)
	if count != 0 {
		t.Error("Image not deleted")
	}
}