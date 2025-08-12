package main

import (
	// The ONLY db import you need now
	"ape-go-services/pkg/db"

	// Internal packages
	v1 "ape-go-services/boilerplate-service/internal/api/v1"
	internalDB "ape-go-services/boilerplate-service/internal/db"

	// Standard and third-party
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system env vars")
	}

	// Connect to PostgreSQL (no changes here)
	if _, err := db.Connect(); err != nil {
		log.Fatalf("Fatal Error: Failed to connect to PostgreSQL: %v", err)
	}
	internalDB.Init()

	// Connect to MongoDB (using renamed functions)
	mongoConfig := db.LoadMongoConfigFromEnv()             // <-- UPDATED
	if err := db.ConnectMongoDB(mongoConfig); err != nil { // <-- UPDATED
		log.Fatalf("Fatal Error: Failed to connect to MongoDB: %v", err)
	}

	setupGracefulShutdown()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	// Health check route (using renamed function)
	router.GET("/health", func(c *gin.Context) {
		var isPostgresConnected bool
		if rawDB, err := db.GetDB().DB(); err == nil && rawDB.Ping() == nil {
			isPostgresConnected = true
		}
		c.JSON(http.StatusOK, gin.H{
			"status":               "ok",
			"mongodb_connected":    db.IsMongoDBConnected(), // <-- UPDATED
			"postgresql_connected": isPostgresConnected,
		})
	})

	v1.SetupRoutes(router)

	log.Printf("Starting boilerplate-service on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutdown signal received. Shutting down gracefully...")
		if err := db.Disconnect(); err != nil { // <-- For Postgres
			log.Printf("Error disconnecting from PostgreSQL: %v", err)
		}
		if err := db.DisconnectMongoDB(); err != nil { // <-- UPDATED for Mongo
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("All connections closed. Exiting application.")
		os.Exit(0)
	}()
}
