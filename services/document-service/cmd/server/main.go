package main

import (
	  "fmt"
	  "log"
	  "os"

	  "github.com/gin-gonic/gin"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    // Healthcheck route
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    log.Printf("Starting document-service on port %s...\n", port)
    err := router.Run(":" + port)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
