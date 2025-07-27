package main

import (
	"log"
	"os"

	"notification-service/internal/config"
	"notification-service/internal/database"
	"notification-service/internal/handlers"
	"notification-service/internal/services"
	"notification-service/internal/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	notificationService := services.NewNotificationService(db)
	schedulerService := scheduler.NewScheduler(notificationService)

	// Start the scheduler
	schedulerService.Start()

	// Initialize handlers
	handler := handlers.NewHandler(notificationService, schedulerService)

	// Setup router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Notification routes
		api.POST("/notifications", handler.SendNotification)
		api.POST("/notifications/schedule", handler.ScheduleNotification)
		api.GET("/notifications", handler.GetNotifications)
		api.GET("/notifications/:id", handler.GetNotification)
		api.PUT("/notifications/:id", handler.UpdateNotification)
		api.DELETE("/notifications/:id", handler.DeleteNotification)

		// Template routes
		api.POST("/templates", handler.CreateTemplate)
		api.GET("/templates", handler.GetTemplates)
		api.GET("/templates/:id", handler.GetTemplate)
		api.PUT("/templates/:id", handler.UpdateTemplate)
		api.DELETE("/templates/:id", handler.DeleteTemplate)

		// Channel routes
		api.GET("/channels", handler.GetChannels)
		api.POST("/channels/test", handler.TestChannel)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "notification-service"})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting notification service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 