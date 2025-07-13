package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bgoettsch/imgonna/backend/internal/handlers"
	"github.com/bgoettsch/imgonna/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configure slog based on environment
	var logLevel slog.Level
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
		logLevel = slog.LevelInfo
	} else {
		// Development mode - enable debug logging
		logLevel = slog.LevelDebug
	}

	// Set up slog with appropriate level
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Starting imgonna backend", 
		"environment", os.Getenv("ENVIRONMENT"),
		"log_level", logLevel.String())

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize services
	anthropicService := services.NewAnthropicService()
	goalsHandler := handlers.NewGoalsHandler(anthropicService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "imgonna-api",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "imgonna API v1",
			})
		})
		
		// Goals endpoint
		api.POST("/goals", goalsHandler.CreateGoal)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server starting", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
		log.Fatal(err)
	}
}