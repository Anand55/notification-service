package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"notification-service/internal/config"
	"notification-service/internal/database"
	"notification-service/internal/models"
	"notification-service/internal/services"
)

// Example usage of the notification service
func main() {
	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize notification service
	notificationService := services.NewNotificationService(db)

	// Example 1: Send an immediate email notification
	fmt.Println("=== Example 1: Send Immediate Email ===")
	emailReq := &models.NotificationRequest{
		Type:       models.EmailNotification,
		Title:      "Welcome to our platform!",
		Message:    "Hello! Welcome to our amazing platform. We're excited to have you on board.",
		Recipient:  "user@example.com",
		Channel:    "default_email",
		Metadata:   models.JSON{"priority": "high"},
	}

	notification, err := notificationService.SendNotification(emailReq)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	} else {
		fmt.Printf("Email notification sent successfully! ID: %d\n", notification.ID)
	}

	// Example 2: Schedule a Slack notification
	fmt.Println("\n=== Example 2: Schedule Slack Notification ===")
	scheduledTime := time.Now().Add(2 * time.Minute) // Schedule for 2 minutes from now
	slackReq := &models.ScheduleRequest{
		NotificationRequest: models.NotificationRequest{
			Type:      models.SlackNotification,
			Title:     "Scheduled Alert",
			Message:   "🚨 This is a scheduled alert notification!",
			Recipient: "#alerts",
			Metadata:  models.JSON{"severity": "medium"},
		},
		ScheduledAt: scheduledTime,
	}

	scheduledNotification, err := notificationService.ScheduleNotification(slackReq)
	if err != nil {
		log.Printf("Failed to schedule Slack notification: %v", err)
	} else {
		fmt.Printf("Slack notification scheduled successfully! ID: %d, Scheduled for: %s\n", 
			scheduledNotification.ID, scheduledNotification.ScheduledAt.Format(time.RFC3339))
	}

	// Example 3: Send in-app notification
	fmt.Println("\n=== Example 3: Send In-App Notification ===")
	inAppReq := &models.NotificationRequest{
		Type:      models.InAppNotification,
		Title:     "New Feature Available",
		Message:   "Check out our new dashboard feature!",
		Recipient: "user123",
		Metadata:  models.JSON{"action_url": "https://app.example.com/dashboard"},
	}

	inAppNotification, err := notificationService.SendNotification(inAppReq)
	if err != nil {
		log.Printf("Failed to send in-app notification: %v", err)
	} else {
		fmt.Printf("In-app notification sent successfully! ID: %d\n", inAppNotification.ID)
	}

	// Example 4: Get notifications with filtering
	fmt.Println("\n=== Example 4: Get Notifications ===")
	notifications, total, err := notificationService.GetNotifications(10, 0, "", "")
	if err != nil {
		log.Printf("Failed to get notifications: %v", err)
	} else {
		fmt.Printf("Found %d notifications (total: %d)\n", len(notifications), total)
		for _, n := range notifications {
			fmt.Printf("- ID: %d, Type: %s, Status: %s, Recipient: %s\n", 
				n.ID, n.Type, n.Status, n.Recipient)
		}
	}

	// Example 5: Create a custom template
	fmt.Println("\n=== Example 5: Create Custom Template ===")
	template := &models.Template{
		Name:    "order_confirmation",
		Type:    models.EmailNotification,
		Subject: "Order Confirmation - #{{.OrderID}}",
		Content: `Hello {{.CustomerName}},

Thank you for your order! Here are your order details:

Order ID: {{.OrderID}}
Total Amount: ${{.Amount}}
Expected Delivery: {{.DeliveryDate}}

You can track your order here: {{.TrackingURL}}

Best regards,
The Team`,
		Variables: models.JSON{
			"CustomerName": "string",
			"OrderID":      "string",
			"Amount":       "string",
			"DeliveryDate": "string",
			"TrackingURL":  "string",
		},
		IsActive: true,
	}

	if err := db.Create(template).Error; err != nil {
		log.Printf("Failed to create template: %v", err)
	} else {
		fmt.Printf("Template created successfully! ID: %d\n", template.ID)
	}

	// Example 6: Use template with data
	fmt.Println("\n=== Example 6: Use Template with Data ===")
	templateReq := &models.NotificationRequest{
		Type:       models.EmailNotification,
		Title:      "Order Confirmation",
		Message:    "Order confirmation message",
		Recipient:  "customer@example.com",
		TemplateID: &template.ID,
		TemplateData: models.JSON{
			"CustomerName": "John Doe",
			"OrderID":      "ORD-12345",
			"Amount":       "99.99",
			"DeliveryDate": "2024-01-20",
			"TrackingURL":  "https://tracking.example.com/ORD-12345",
		},
	}

	templatedNotification, err := notificationService.SendNotification(templateReq)
	if err != nil {
		log.Printf("Failed to send templated notification: %v", err)
	} else {
		fmt.Printf("Templated notification sent successfully! ID: %d\n", templatedNotification.ID)
		fmt.Printf("Processed message: %s\n", templatedNotification.Message)
	}

	// Example 7: Test channel connections
	fmt.Println("\n=== Example 7: Test Channel Connections ===")
	
	// Test email connection
	if err := notificationService.GetEmailSender().TestConnection(); err != nil {
		fmt.Printf("Email connection test failed: %v\n", err)
	} else {
		fmt.Println("Email connection test successful!")
	}

	// Test Slack connection
	if err := notificationService.GetSlackSender().TestConnection(); err != nil {
		fmt.Printf("Slack connection test failed: %v\n", err)
	} else {
		fmt.Println("Slack connection test successful!")
	}

	// Test in-app connection
	if err := notificationService.GetInAppSender().TestConnection(); err != nil {
		fmt.Printf("In-app connection test failed: %v\n", err)
	} else {
		fmt.Println("In-app connection test successful!")
	}

	fmt.Println("\n=== Examples completed! ===")
}

// Helper function to pretty print JSON
func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(b))
} 