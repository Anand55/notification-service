package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"notification-service/internal/config"
	"notification-service/internal/models"

	"gorm.io/gorm"
)

// NotificationService handles notification operations
type NotificationService struct {
	db     *gorm.DB
	config *config.Config
	emailSender *EmailSender
	slackSender *SlackSender
	inAppSender *InAppSender
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *gorm.DB) *NotificationService {
	cfg := config.Load()
	
	return &NotificationService{
		db:     db,
		config: cfg,
		emailSender: NewEmailSender(cfg),
		slackSender: NewSlackSender(cfg),
		inAppSender: NewInAppSender(),
	}
}

// SendNotification sends a notification immediately
func (s *NotificationService) SendNotification(req *models.NotificationRequest) (*models.Notification, error) {
	// Create notification record
	notification := &models.Notification{
		Type:       req.Type,
		Status:     models.PendingStatus,
		Title:      req.Title,
		Message:    req.Message,
		Recipient:  req.Recipient,
		Channel:    req.Channel,
		TemplateID: req.TemplateID,
		Metadata:   req.Metadata,
	}

	// Process template if provided
	if req.TemplateID != nil {
		if err := s.processTemplate(notification, req.TemplateData); err != nil {
			return nil, err
		}
	}

	// Save to database
	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	// Send notification
	if err := s.sendNotification(notification); err != nil {
		notification.Status = models.FailedStatus
		s.db.Save(notification)
		return notification, err
	}

	// Update status to sent
	notification.Status = models.SentStatus
	now := time.Now()
	notification.SentAt = &now
	s.db.Save(notification)

	return notification, nil
}

// ScheduleNotification schedules a notification for later
func (s *NotificationService) ScheduleNotification(req *models.ScheduleRequest) (*models.Notification, error) {
	notification := &models.Notification{
		Type:        req.Type,
		Status:      models.ScheduledStatus,
		Title:       req.Title,
		Message:     req.Message,
		Recipient:   req.Recipient,
		Channel:     req.Channel,
		TemplateID:  req.TemplateID,
		Metadata:    req.Metadata,
		ScheduledAt: &req.ScheduledAt,
	}

	// Process template if provided
	if req.TemplateID != nil {
		if err := s.processTemplate(notification, req.TemplateData); err != nil {
			return nil, err
		}
	}

	// Save to database
	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

// ProcessScheduledNotifications processes all scheduled notifications that are due
func (s *NotificationService) ProcessScheduledNotifications() error {
	var notifications []models.Notification
	
	if err := s.db.Where("status = ? AND scheduled_at <= ?", 
		models.ScheduledStatus, time.Now()).Find(&notifications).Error; err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := s.sendNotification(&notification); err != nil {
			log.Printf("Failed to send scheduled notification %d: %v", notification.ID, err)
			notification.Status = models.FailedStatus
			s.db.Save(&notification)
			continue
		}

		notification.Status = models.SentStatus
		now := time.Now()
		notification.SentAt = &now
		s.db.Save(&notification)
	}

	return nil
}

// GetNotifications retrieves notifications with optional filtering
func (s *NotificationService) GetNotifications(limit, offset int, status models.NotificationStatus, notificationType models.NotificationType) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Preload("Template").Limit(limit).Offset(offset).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetNotification retrieves a single notification by ID
func (s *NotificationService) GetNotification(id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := s.db.Preload("Template").First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// UpdateNotification updates a notification
func (s *NotificationService) UpdateNotification(id uint, updates map[string]interface{}) (*models.Notification, error) {
	var notification models.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&notification).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(id uint) error {
	return s.db.Delete(&models.Notification{}, id).Error
}

// GetDB returns the database instance
func (s *NotificationService) GetDB() *gorm.DB {
	return s.db
}

// GetEmailSender returns the email sender
func (s *NotificationService) GetEmailSender() *EmailSender {
	return s.emailSender
}

// GetSlackSender returns the Slack sender
func (s *NotificationService) GetSlackSender() *SlackSender {
	return s.slackSender
}

// GetInAppSender returns the in-app sender
func (s *NotificationService) GetInAppSender() *InAppSender {
	return s.inAppSender
}

// sendNotification sends a notification through the appropriate channel
func (s *NotificationService) sendNotification(notification *models.Notification) error {
	switch notification.Type {
	case models.EmailNotification:
		return s.emailSender.Send(notification)
	case models.SlackNotification:
		return s.slackSender.Send(notification)
	case models.InAppNotification:
		return s.inAppSender.Send(notification)
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}

// processTemplate processes a template with the provided data
func (s *NotificationService) processTemplate(notification *models.Notification, templateData models.JSON) error {
	if notification.TemplateID == nil {
		return nil
	}

	var tmpl models.Template
	if err := s.db.First(&tmpl, *notification.TemplateID).Error; err != nil {
		return err
	}

	// Parse template
	t, err := template.New("notification").Parse(tmpl.Content)
	if err != nil {
		return err
	}

	// Execute template
	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData); err != nil {
		return err
	}

	// Update notification with processed content
	notification.Message = buf.String()
	if tmpl.Subject != "" {
		notification.Title = tmpl.Subject
	}

	return nil
} 