package main

import (
	  "log"
	  "os"

	  "github.com/gin-gonic/gin"
	  "github.com/joho/godotenv"
)

func main() {

    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, continuing with system env vars")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    err := router.SetTrustedProxies(nil)
    if err != nil {
        log.Fatalf("Error setting trusted proxies: %v", err)
    }

    // Healthcheck route
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    log.Printf("Starting user-cv-service on port %s...\n", port)
    err = router.Run(":" + port)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
