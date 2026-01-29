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

func (s *UserService) Register(user *User) (*User, error) {
	user.Password = hashPassword(user.Password)
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
