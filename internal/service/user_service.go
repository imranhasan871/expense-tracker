package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"time"
)

type UserService interface {
	CreateUser(name, displayID, email string, role models.UserRole) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUserRole(userID int, role models.UserRole) error
}

type userService struct {
	userRepo     repository.UserRepository
	emailService EmailService
}

func NewUserService(userRepo repository.UserRepository, emailService EmailService) UserService {
	return &userService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *userService) CreateUser(name, displayID, email string, role models.UserRole) (*models.User, error) {
	if name == "" || displayID == "" || email == "" {
		return nil, errors.New("name, user ID, and email are required")
	}

	token, _ := generateRandomToken(32)
	expiry := time.Now().Add(24 * time.Hour)

	user := &models.User{
		Username:               name,
		UserDisplayID:          displayID,
		Email:                  email,
		Role:                   role,
		IsActive:               false,
		PasswordSetToken:       &token,
		PasswordSetTokenExpiry: &expiry,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Send automated email
	if err := s.emailService.SendPasswordSetEmail(email, token); err != nil {
		// We might want to log this but not fail the whole process
		// For now, we'll return the error
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateUserRole(userID int, role models.UserRole) error {
	// Validate role
	if role != models.RoleAdmin && role != models.RoleManagement && role != models.RoleExecutive {
		return errors.New("invalid role")
	}

	return s.userRepo.UpdateRole(userID, role)
}

func generateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
