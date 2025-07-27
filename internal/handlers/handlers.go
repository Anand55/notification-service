package handlers

import (
	"net/http"
	"strconv"

	"notification-service/internal/models"
	"notification-service/internal/scheduler"
	"notification-service/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler handles HTTP requests
type Handler struct {
	notificationService *services.NotificationService
	schedulerService    *scheduler.Scheduler
}

// NewHandler creates a new handler
func NewHandler(notificationService *services.NotificationService, schedulerService *scheduler.Scheduler) *Handler {
	return &Handler{
		notificationService: notificationService,
		schedulerService:    schedulerService,
	}
}

// SendNotification handles sending immediate notifications
func (h *Handler) SendNotification(c *gin.Context) {
	var req models.NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationService.SendNotification(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// ScheduleNotification handles scheduling notifications
func (h *Handler) ScheduleNotification(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationService.ScheduleNotification(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetNotifications handles retrieving notifications with pagination and filtering
func (h *Handler) GetNotifications(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := models.NotificationStatus(c.Query("status"))
	notificationType := models.NotificationType(c.Query("type"))

	notifications, total, err := h.notificationService.GetNotifications(limit, offset, status, notificationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	})
}

// GetNotification handles retrieving a single notification
func (h *Handler) GetNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	notification, err := h.notificationService.GetNotification(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// UpdateNotification handles updating a notification
func (h *Handler) UpdateNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationService.UpdateNotification(uint(id), updates)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// DeleteNotification handles deleting a notification
func (h *Handler) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := h.notificationService.DeleteNotification(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

// CreateTemplate handles creating a new template
func (h *Handler) CreateTemplate(c *gin.Context) {
	var req models.TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := &models.Template{
		Name:      req.Name,
		Type:      req.Type,
		Subject:   req.Subject,
		Content:   req.Content,
		Variables: req.Variables,
		IsActive:  true,
	}

	if err := h.notificationService.GetDB().Create(template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// GetTemplates handles retrieving templates
func (h *Handler) GetTemplates(c *gin.Context) {
	var templates []models.Template
	if err := h.notificationService.GetDB().Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetTemplate handles retrieving a single template
func (h *Handler) GetTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var template models.Template
	if err := h.notificationService.GetDB().First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate handles updating a template
func (h *Handler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var req models.TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template models.Template
	if err := h.notificationService.GetDB().First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	template.Name = req.Name
	template.Type = req.Type
	template.Subject = req.Subject
	template.Content = req.Content
	template.Variables = req.Variables

	if err := h.notificationService.GetDB().Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate handles deleting a template
func (h *Handler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	if err := h.notificationService.GetDB().Delete(&models.Template{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// GetChannels handles retrieving available channels
func (h *Handler) GetChannels(c *gin.Context) {
	var channels []models.Channel
	if err := h.notificationService.GetDB().Find(&channels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channels)
}

// TestChannel handles testing a notification channel
func (h *Handler) TestChannel(c *gin.Context) {
	var req struct {
		Type string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch models.NotificationType(req.Type) {
	case models.EmailNotification:
		err = h.notificationService.GetEmailSender().TestConnection()
	case models.SlackNotification:
		err = h.notificationService.GetSlackSender().TestConnection()
	case models.InAppNotification:
		err = h.notificationService.GetInAppSender().TestConnection()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported notification type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel test successful"})
} 