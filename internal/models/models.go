package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	EmailNotification    NotificationType = "email"
	SlackNotification    NotificationType = "slack"
	InAppNotification    NotificationType = "in_app"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	PendingStatus   NotificationStatus = "pending"
	SentStatus      NotificationStatus = "sent"
	FailedStatus    NotificationStatus = "failed"
	ScheduledStatus NotificationStatus = "scheduled"
)

// Notification represents a notification record
type Notification struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	Type        NotificationType   `json:"type" gorm:"not null"`
	Status      NotificationStatus `json:"status" gorm:"not null;default:'pending'"`
	Title       string             `json:"title" gorm:"not null"`
	Message     string             `json:"message" gorm:"not null"`
	Recipient   string             `json:"recipient" gorm:"not null"`
	Channel     string             `json:"channel"`
	TemplateID  *uint              `json:"template_id"`
	Template    *Template          `json:"template,omitempty"`
	ScheduledAt *time.Time         `json:"scheduled_at"`
	SentAt      *time.Time         `json:"sent_at"`
	Metadata    JSON               `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"deleted_at,omitempty" gorm:"index"`
}

// Template represents a notification template
type Template struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;unique"`
	Type        NotificationType `json:"type" gorm:"not null"`
	Subject     string         `json:"subject"`
	Content     string         `json:"content" gorm:"not null"`
	Variables   JSON           `json:"variables" gorm:"type:json"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Channel represents a notification channel configuration
type Channel struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;unique"`
	Type        NotificationType `json:"type" gorm:"not null"`
	Config      JSON           `json:"config" gorm:"type:json"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// JSON is a custom type for JSON fields
type JSON map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, j)
}

// NotificationRequest represents the request structure for sending notifications
type NotificationRequest struct {
	Type        NotificationType `json:"type" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Message     string           `json:"message" binding:"required"`
	Recipient   string           `json:"recipient" binding:"required"`
	Channel     string           `json:"channel"`
	TemplateID  *uint            `json:"template_id"`
	TemplateData JSON            `json:"template_data"`
	Metadata    JSON             `json:"metadata"`
}

// ScheduleRequest represents the request structure for scheduling notifications
type ScheduleRequest struct {
	NotificationRequest
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

// TemplateRequest represents the request structure for templates
type TemplateRequest struct {
	Name      string           `json:"name" binding:"required"`
	Type      NotificationType `json:"type" binding:"required"`
	Subject   string           `json:"subject"`
	Content   string           `json:"content" binding:"required"`
	Variables JSON             `json:"variables"`
}

// ChannelRequest represents the request structure for channels
type ChannelRequest struct {
	Name   string           `json:"name" binding:"required"`
	Type   NotificationType `json:"type" binding:"required"`
	Config JSON             `json:"config" binding:"required"`
} 