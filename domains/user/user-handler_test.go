package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vixel/config"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestUserHandler(t *testing.T) (*UserHandler, *gorm.DB) {
	config.Config.JWTSecret = "test_secret"
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	db.AutoMigrate(&User{})
	service := NewUserService(db)
	handler := NewUserHandler(service)
	return handler, db
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupTestUserHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	dto := RegisterDto{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(dto)

	c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register()(c)

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

	token, ok := data["access_token"].(string)
	if !ok || token == "" {
		t.Error("Access token should not be empty")
	}
}

func TestUserHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, db := setupTestUserHandler(t)

	// Pre-register user
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	service := NewUserService(db)
	_, err := service.Register(user)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	dto := LoginDto{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(dto)

	c.Request = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login()(c)

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

	token, ok := data["access_token"].(string)
	if !ok || token == "" {
		t.Error("Access token should not be empty")
	}
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupTestUserHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	dto := LoginDto{
		Email:    "invalid@example.com",
		Password: "password",
	}
	jsonData, _ := json.Marshal(dto)

	c.Request = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login()(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
