package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"gorm.io/gorm"

	"github.com/flux-panel/go-backend/models"
)

// UserService contains business logic for users.
type UserService struct {
	DB *gorm.DB
}

// NewUserService creates a new UserService.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// hashPassword returns MD5 hash of password.
func hashPassword(pwd string) string {
	h := md5.Sum([]byte(pwd))
	return hex.EncodeToString(h[:])
}

// CreateUser inserts a new user if username is unique.
func (s *UserService) CreateUser(u *models.User) error {
	var count int64
	if err := s.DB.Model(&models.User{}).Where("user = ?", u.User).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}
	u.Password = hashPassword(u.Password)
	return s.DB.Create(u).Error
}

// Login validates credentials and returns user.
func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("user = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	if user.Password != hashPassword(password) {
		return nil, errors.New("invalid credentials")
	}
	if user.Status != 1 {
		return nil, errors.New("account disabled")
	}
	return &user, nil
}

// ListUsers returns all non-admin users.
func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := s.DB.Where("role_id <> ?", 0).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
