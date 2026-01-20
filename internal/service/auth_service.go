package service

import (
	"errors"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"

	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (*models.User, string, error)
	SetPassword(token, password string) error
	ValidateToken(token string) (*models.User, error)
	Logout(sessionToken string) error
	IsAuthenticated(sessionToken string) (*models.User, bool)
}

type authService struct {
	userRepo     repository.UserRepository
	sessionStore map[string]*models.User
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo:     userRepo,
		sessionStore: make(map[string]*models.User),
	}
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	email = strings.ToLower(email)
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, "", errors.New("account is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	sessionToken, _ := generateRandomToken(32)
	s.sessionStore[sessionToken] = user

	return user, sessionToken, nil
}

func (s *authService) SetPassword(token, password string) error {
	user, err := s.userRepo.GetByToken(token)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(user.ID, string(hashedPassword))
}

func (s *authService) ValidateToken(token string) (*models.User, error) {
	return s.userRepo.GetByToken(token)
}

func (s *authService) Logout(sessionToken string) error {
	delete(s.sessionStore, sessionToken)
	return nil
}

func (s *authService) IsAuthenticated(sessionToken string) (*models.User, bool) {
	user, ok := s.sessionStore[sessionToken]
	return user, ok
}
