package main

import (
	// Shared packages from the monorepo root
	"ape-go-services/pkg/db"
	"ape-go-services/pkg/mongodb"

	// Internal packages for this specific service
	v1 "ape-go-services/website-cms-service/internal/api/v1"
	internalDB "ape-go-services/website-cms-service/internal/db" // Aliased import to avoid name conflict

	// Standard library imports
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// Third-party imports
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system env vars")
	}

	// ---- Database Connections ----

	// 1. Connect to PostgreSQL using the shared package
	if _, err := db.Connect(); err != nil {
		log.Fatalf("Fatal Error: Failed to connect to PostgreSQL: %v", err)
	}

	// 2. Initialize this service's internal DB handle with the established connection
	internalDB.Init()

	// 3. Connect to MongoDB
	mongoConfig := mongodb.LoadConfigFromEnv()
	if err := mongodb.Connect(mongoConfig); err != nil {
		log.Fatalf("Fatal Error: Failed to connect to MongoDB: %v", err)
	}

	// ---- Graceful Shutdown ----
	// This must be set up after connections are established.
	setupGracefulShutdown()

	// ---- Gin Server Setup ----
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	// ---- Routes ----

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		// Ping internal services to check their status
		mongoConnected := mongodb.IsConnected()

		// For postgres, we can get the rawDB and ping it
		var isPostgresConnected bool
		if rawDB, err := db.GetDB().DB(); err == nil {
			if err := rawDB.Ping(); err == nil {
				isPostgresConnected = true
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":               "ok",
			"mongodb_connected":    mongoConnected,
			"postgresql_connected": isPostgresConnected,
		})
	})

	// Setup API v1 routes
	v1.SetupRoutes(router)

	// ---- Start Server ----
	log.Printf("Starting website-cms-service on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupGracefulShutdown handles termination signals to ensure all connections are closed properly.
func setupGracefulShutdown() {
	// Create a channel to listen for OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start a new goroutine to block and wait for a signal
	go func() {
		<-c // This blocks until a signal is received

		log.Println("Shutdown signal received. Shutting down gracefully...")

		// Disconnect from PostgreSQL
		if err := db.Disconnect(); err != nil {
			log.Printf("Error disconnecting from PostgreSQL: %v", err)
		}

		// Disconnect from MongoDB
		if err := mongodb.Disconnect(); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}

		log.Println("All connections closed. Exiting application.")
		os.Exit(0)
	}()
}
