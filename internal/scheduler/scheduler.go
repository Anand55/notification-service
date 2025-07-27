package scheduler

import (
	"log"
	"time"

	"notification-service/internal/services"

	"github.com/go-co-op/gocron"
)

// Scheduler handles scheduled notifications
type Scheduler struct {
	scheduler *gocron.Scheduler
	notificationService *services.NotificationService
}

// NewScheduler creates a new scheduler
func NewScheduler(notificationService *services.NotificationService) *Scheduler {
	return &Scheduler{
		scheduler: gocron.NewScheduler(time.UTC),
		notificationService: notificationService,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	// Process scheduled notifications every minute
	s.scheduler.Every(1).Minute().Do(s.processScheduledNotifications)
	
	// Start the scheduler
	s.scheduler.StartAsync()
	
	log.Println("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.scheduler.Stop()
	log.Println("Scheduler stopped")
}

// processScheduledNotifications processes all scheduled notifications that are due
func (s *Scheduler) processScheduledNotifications() {
	if err := s.notificationService.ProcessScheduledNotifications(); err != nil {
		log.Printf("Error processing scheduled notifications: %v", err)
	}
}

// GetScheduler returns the underlying gocron scheduler
func (s *Scheduler) GetScheduler() *gocron.Scheduler {
	return s.scheduler
} 