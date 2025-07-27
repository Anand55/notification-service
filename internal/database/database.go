package database

import (
	"notification-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes the database connection and runs migrations
func Init(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.AutoMigrate(
		&models.Notification{},
		&models.Template{},
		&models.Channel{},
	); err != nil {
		return nil, err
	}

	// Seed default data
	if err := seedDefaultData(db); err != nil {
		return nil, err
	}

	return db, nil
}

// seedDefaultData seeds the database with default templates and channels
func seedDefaultData(db *gorm.DB) error {
	// Check if default templates already exist
	var count int64
	db.Model(&models.Template{}).Count(&count)
	if count > 0 {
		return nil // Data already seeded
	}

	// Create default email templates
	defaultTemplates := []models.Template{
		{
			Name:    "welcome_email",
			Type:    models.EmailNotification,
			Subject: "Welcome to our platform!",
			Content: "Hello {{.Name}},\n\nWelcome to our platform! We're excited to have you on board.\n\nBest regards,\nThe Team",
			Variables: models.JSON{
				"Name": "string",
			},
			IsActive: true,
		},
		{
			Name:    "password_reset",
			Type:    models.EmailNotification,
			Subject: "Password Reset Request",
			Content: "Hello {{.Name}},\n\nYou requested a password reset. Click the link below to reset your password:\n\n{{.ResetLink}}\n\nIf you didn't request this, please ignore this email.\n\nBest regards,\nThe Team",
			Variables: models.JSON{
				"Name":      "string",
				"ResetLink": "string",
			},
			IsActive: true,
		},
		{
			Name:    "slack_alert",
			Type:    models.SlackNotification,
			Content: "🚨 Alert: {{.AlertType}}\n\n{{.Message}}\n\nTime: {{.Timestamp}}\nSeverity: {{.Severity}}",
			Variables: models.JSON{
				"AlertType": "string",
				"Message":   "string",
				"Timestamp": "string",
				"Severity":  "string",
			},
			IsActive: true,
		},
		{
			Name:    "in_app_notification",
			Type:    models.InAppNotification,
			Content: "{{.Title}}\n\n{{.Message}}\n\n{{.ActionText}}: {{.ActionUrl}}",
			Variables: models.JSON{
				"Title":     "string",
				"Message":   "string",
				"ActionText": "string",
				"ActionUrl": "string",
			},
			IsActive: true,
		},
	}

	for _, template := range defaultTemplates {
		if err := db.Create(&template).Error; err != nil {
			return err
		}
	}

	// Create default channels
	defaultChannels := []models.Channel{
		{
			Name:   "default_email",
			Type:   models.EmailNotification,
			Config: models.JSON{},
			IsActive: true,
		},
		{
			Name:   "default_slack",
			Type:   models.SlackNotification,
			Config: models.JSON{},
			IsActive: true,
		},
		{
			Name:   "default_in_app",
			Type:   models.InAppNotification,
			Config: models.JSON{},
			IsActive: true,
		},
	}

	for _, channel := range defaultChannels {
		if err := db.Create(&channel).Error; err != nil {
			return err
		}
	}

	return nil
} 