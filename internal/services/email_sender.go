package services

import (
	"fmt"
	"notification-service/internal/config"
	"notification-service/internal/models"

	"gopkg.in/gomail.v2"
)

// EmailSender handles email notifications
type EmailSender struct {
	config *config.Config
}

// NewEmailSender creates a new email sender
func NewEmailSender(config *config.Config) *EmailSender {
	return &EmailSender{
		config: config,
	}
}

// Send sends an email notification
func (e *EmailSender) Send(notification *models.Notification) error {
	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.EmailUsername)
	m.SetHeader("To", notification.Recipient)
	m.SetHeader("Subject", notification.Title)
	m.SetBody("text/plain", notification.Message)

	// Add HTML body if metadata contains HTML content
	if htmlContent, exists := notification.Metadata["html_content"]; exists {
		if htmlStr, ok := htmlContent.(string); ok {
			m.SetBody("text/html", htmlStr)
		}
	}

	// Create dialer
	d := gomail.NewDialer(
		e.config.EmailHost,
		e.config.EmailPort,
		e.config.EmailUsername,
		e.config.EmailPassword,
	)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// TestConnection tests the email connection
func (e *EmailSender) TestConnection() error {
	d := gomail.NewDialer(
		e.config.EmailHost,
		e.config.EmailPort,
		e.config.EmailUsername,
		e.config.EmailPassword,
	)

	s, err := d.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer s.Close()

	return nil
} 