package main

import (
	"ape-go-services/pkg/mongodb"
	v1 "ape-go-services/website-cms-service/internal/api/v1"
	"ape-go-services/website-cms-service/internal/db"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system env vars")
	}

	// Connect to PostgreSQL (existing)
	db.Connect()

	// Connect to MongoDB (new)
	mongoConfig := mongodb.LoadConfigFromEnv()
	if err := mongodb.Connect(mongoConfig); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		mongoConnected := mongodb.IsConnected()
		c.JSON(200, gin.H{
			"status":            "ok",
			"mongodb_connected": mongoConnected,
		})
	})

	// Setup API routes
	v1.SetupRoutes(router)

	log.Printf("Starting website-cms-service on port %s...", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")

		// Disconnect from MongoDB
		if err := mongodb.Disconnect(); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}

		os.Exit(0)
	}()
}
