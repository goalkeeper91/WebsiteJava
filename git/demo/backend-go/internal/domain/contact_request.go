package domain

import (
	"time"
)

type ContactRequest struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Subject      string    `json:"subject"`
	Message      string    `json:"message"`
	ConsentGiven bool      `json:"consent_given"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewContactRequest(name, email, phone, subject, message string, consentGiven bool) *ContactRequest {
	return &ContactRequest{
		Name:         name,
		Email:        email,
		Phone:        phone,
		Subject:      subject,
		Message:      message,
		ConsentGiven: consentGiven,
		CreatedAt:    time.Now(),
	}
}

func (c *ContactRequest) Validate() error {
	if c.Name == "" {
		return ErrInvalidInput
	}
	if c.Email == "" {
		return ErrInvalidInput
	}
	if c.Subject == "" {
		return ErrInvalidInput
	}
	if c.Message == "" {
		return ErrInvalidInput
	}
	if !c.ConsentGiven {
		return ErrInvalidInput
	}
	return nil
}