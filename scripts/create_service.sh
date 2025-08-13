#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "❌ Error: you must provide a service name (e.g., user-service)"
  exit 1
fi

SERVICE_NAME=$1
BASE_MODULE="ape-go-services"
SERVICES_DIR="services"
SERVICE_PATH="$SERVICES_DIR/$SERVICE_NAME"

if [ -d "$SERVICE_PATH" ]; then
  echo "⚠️ Service '$SERVICE_NAME' already exists at $SERVICE_PATH"
  exit 1
fi

echo "🚀 Creating service: $SERVICE_NAME"

# Create folder structure
mkdir -p $SERVICE_PATH
mkdir -p $SERVICE_PATH/cmd/server
mkdir -p $SERVICE_PATH/internal/api/v1/{dto,handlers,middleware}
mkdir -p $SERVICE_PATH/internal/{auth,config,db,mappers,models,repositories,services}
mkdir -p $SERVICE_PATH/docs

# Add .gitkeep files to preserve structure
find $SERVICE_PATH/internal -type d -empty -exec touch {}/.gitkeep \;
touch $SERVICE_PATH/docs/.gitkeep

# main.go with basic Gin example
cat <<EOF > $SERVICE_PATH/cmd/server/main.go
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

    log.Printf("Starting $SERVICE_NAME on port %s...\n", port)
    err = router.Run(":" + port)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
EOF

# go.mod
cat <<EOF > $SERVICE_PATH/go.mod
module $BASE_MODULE/$SERVICE_NAME

go 1.24
EOF

# run go mod tidy inside the service
(
	cd $SERVICE_PATH
	go mod tidy
)

# routes.go placeholder
cat <<EOF > $SERVICE_PATH/internal/api/v1/routes.go
package v1

// Define your routes and route groups here
EOF

# Dockerfile
cat <<EOF > $SERVICE_PATH/Dockerfile
# Build stage
FROM golang:1.24 AS builder

# Set working dir before copying
WORKDIR /app

# Copy the entire monorepo (from root)
COPY . .

# Move into the service directory
WORKDIR /app/services/$SERVICE_NAME

RUN go mod tidy && go build -o main ./cmd/server

# Final stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/services/$SERVICE_NAME/main .
COPY --from=builder /app/services/$SERVICE_NAME/.env .

CMD ["./main"]
EOF

# Makefile
cat <<EOF > $SERVICE_PATH/Makefile
# Build the service binary
build:
	go build -o bin/$SERVICE_NAME ./cmd/server

# Run the service locally
run:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy
EOF

# Add to go.work
if [ -f "go.work" ]; then
	echo "🧩 Adding $SERVICE_NAME to go.work"
	go work use ./$SERVICE_PATH
else
	echo "🧩 Creating go.work"
	go work init ./$SERVICE_PATH
fi

echo "✅ Service '$SERVICE_NAME' created successfully!"
