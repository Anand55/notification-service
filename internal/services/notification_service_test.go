package services

import (
	"testing"
	"time"

	"notification-service/internal/database"
	"notification-service/internal/models"
)

func TestNotificationService(t *testing.T) {
	// Initialize test database
	db, err := database.Init(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Initialize notification service
	service := NewNotificationService(db)

	// Test 1: Send immediate notification
	t.Run("Send Immediate Notification", func(t *testing.T) {
		req := &models.NotificationRequest{
			Type:      models.InAppNotification, // Use in-app to avoid email config issues
			Title:     "Test Notification",
			Message:   "This is a test notification",
			Recipient: "test@example.com",
		}

		notification, err := service.SendNotification(req)
		if err != nil {
			t.Errorf("Failed to send notification: %v", err)
		}

		if notification.ID == 0 {
			t.Error("Notification ID should not be zero")
		}

		if notification.Status != models.SentStatus {
			t.Errorf("Expected status %s, got %s", models.SentStatus, notification.Status)
		}
	})

	// Test 2: Schedule notification
	t.Run("Schedule Notification", func(t *testing.T) {
		scheduledTime := time.Now().Add(1 * time.Hour)
		req := &models.ScheduleRequest{
			NotificationRequest: models.NotificationRequest{
				Type:      models.SlackNotification,
				Title:     "Scheduled Test",
				Message:   "This is a scheduled test",
				Recipient: "#test",
			},
			ScheduledAt: scheduledTime,
		}

		notification, err := service.ScheduleNotification(req)
		if err != nil {
			t.Errorf("Failed to schedule notification: %v", err)
		}

		if notification.ID == 0 {
			t.Error("Notification ID should not be zero")
		}

		if notification.Status != models.ScheduledStatus {
			t.Errorf("Expected status %s, got %s", models.ScheduledStatus, notification.Status)
		}

		if notification.ScheduledAt == nil {
			t.Error("ScheduledAt should not be nil")
		}
	})

	// Test 3: Get notifications
	t.Run("Get Notifications", func(t *testing.T) {
		notifications, total, err := service.GetNotifications(10, 0, "", "")
		if err != nil {
			t.Errorf("Failed to get notifications: %v", err)
		}

		if total < 2 {
			t.Errorf("Expected at least 2 notifications, got %d", total)
		}

		if len(notifications) < 2 {
			t.Errorf("Expected at least 2 notifications in result, got %d", len(notifications))
		}
	})

	// Test 4: Get notification by ID
	t.Run("Get Notification by ID", func(t *testing.T) {
		// First, get a notification
		notifications, _, err := service.GetNotifications(1, 0, "", "")
		if err != nil {
			t.Errorf("Failed to get notifications: %v", err)
		}

		if len(notifications) == 0 {
			t.Skip("No notifications to test with")
		}

		notification, err := service.GetNotification(notifications[0].ID)
		if err != nil {
			t.Errorf("Failed to get notification by ID: %v", err)
		}

		if notification.ID != notifications[0].ID {
			t.Errorf("Expected notification ID %d, got %d", notifications[0].ID, notification.ID)
		}
	})

	// Test 5: Update notification
	t.Run("Update Notification", func(t *testing.T) {
		// First, get a notification
		notifications, _, err := service.GetNotifications(1, 0, "", "")
		if err != nil {
			t.Errorf("Failed to get notifications: %v", err)
		}

		if len(notifications) == 0 {
			t.Skip("No notifications to test with")
		}

		updates := map[string]interface{}{
			"title": "Updated Title",
		}

		notification, err := service.UpdateNotification(notifications[0].ID, updates)
		if err != nil {
			t.Errorf("Failed to update notification: %v", err)
		}

		if notification.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", notification.Title)
		}
	})

	// Test 6: Delete notification
	t.Run("Delete Notification", func(t *testing.T) {
		// First, get a notification
		notifications, _, err := service.GetNotifications(1, 0, "", "")
		if err != nil {
			t.Errorf("Failed to get notifications: %v", err)
		}

		if len(notifications) == 0 {
			t.Skip("No notifications to test with")
		}

		err = service.DeleteNotification(notifications[0].ID)
		if err != nil {
			t.Errorf("Failed to delete notification: %v", err)
		}

		// Verify it's deleted
		_, err = service.GetNotification(notifications[0].ID)
		if err == nil {
			t.Error("Notification should be deleted")
		}
	})
}

func TestTemplateProcessing(t *testing.T) {
	// Initialize test database
	db, err := database.Init(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Initialize notification service
	service := NewNotificationService(db)

	// Create a test template
	template := &models.Template{
		Name:    "test_template",
		Type:    models.EmailNotification,
		Subject: "Hello {{.Name}}",
		Content: "Hello {{.Name}},\n\nWelcome to {{.Platform}}!",
		Variables: models.JSON{
			"Name":     "string",
			"Platform": "string",
		},
		IsActive: true,
	}

	if err := db.Create(template).Error; err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Test template processing
	t.Run("Process Template", func(t *testing.T) {
		req := &models.NotificationRequest{
			Type:       models.InAppNotification, // Use in-app to avoid email config issues
			Title:      "Test",
			Message:    "Test message",
			Recipient:  "test@example.com",
			TemplateID: &template.ID,
			TemplateData: models.JSON{
				"Name":     "John Doe",
				"Platform": "Our Platform",
			},
		}

		notification, err := service.SendNotification(req)
		if err != nil {
			t.Errorf("Failed to send notification with template: %v", err)
		}

		expectedMessage := "Hello John Doe,\n\nWelcome to Our Platform!"
		if notification.Message != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, notification.Message)
		}
	})
} 