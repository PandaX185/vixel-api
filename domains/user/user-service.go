package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func (s *UserService) Register(user *User) (*TokenResponse, error) {
	user.Password = hashPassword(user.Password)
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	token, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken: token,
	}, nil
}

func (s *UserService) Login(email, password string) (*TokenResponse, error) {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	token, err := generateAccessToken(&user)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken: token,
	}, nil
}
