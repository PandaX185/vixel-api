package user

import (
	"testing"

	"vixel/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	config.Config.JWTSecret = "test_secret"
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	db.AutoMigrate(&User{})
	return db
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashed := hashPassword(password)

	if hashed == password {
		t.Error("Password should be hashed")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		t.Errorf("Hashed password should match original: %v", err)
	}
}

func TestUserService_Register(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	token, err := service.Register(user)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}

	var createdUser User
	if err := db.Where("email = ?", "test@example.com").First(&createdUser).Error; err != nil {
		t.Fatalf("User not found in DB: %v", err)
	}

	if createdUser.Username != "testuser" {
		t.Error("Username not saved correctly")
	}

	if createdUser.Password == "password123" {
		t.Error("Password should be hashed")
	}
}

func TestUserService_Login_Success(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := service.Register(user)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	token, err := service.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}
}

func TestUserService_Login_InvalidEmail(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	_, err := service.Login("nonexistent@example.com", "password")
	if err == nil {
		t.Error("Expected error for invalid email")
	}
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := service.Register(user)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	_, err = service.Login("test@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected error for invalid password")
	}
}