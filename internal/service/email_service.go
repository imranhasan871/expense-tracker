package service

import (
	"log"
)

type EmailService interface {
	SendPasswordSetEmail(email, token string) error
}

type mockEmailService struct{}

func NewEmailService() EmailService {
	return &mockEmailService{}
}

func (s *mockEmailService) SendPasswordSetEmail(email, token string) error {
	// In a real app, this would send an actual email.
	// For now, we'll just log it.
	resetLink := "http://localhost:8080/set-password?token=" + token
	log.Printf("[EMAIL MOCK] To: %s | Subject: Set Your Password | Message: Please click the link to set your password: %s", email, resetLink)
	return nil
}
