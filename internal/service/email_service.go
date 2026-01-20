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
	resetLink := "http://localhost:8080/set-password?token=" + token
	log.Printf("[EMAIL MOCK] To: %s | Subject: Set Your Password | Message: Please click the link to set your password: %s", email, resetLink)
	return nil
}
