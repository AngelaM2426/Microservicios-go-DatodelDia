package db

import (
	"ape-go-services/pkg/db" // Your shared database package
	"log"

	"gorm.io/gorm"
)

// DB is a package-level variable that holds the database instance for this service.
var DB *gorm.DB

// Init initializes the database connection for the service.
// It retrieves the connection established in main.go from the shared db package.
func Init() {
	// Get the singleton instance from the shared package.
	// This ensures we are using the same connection that was
	// initialized and tested at startup.
	DB = db.GetDB()
	if DB == nil {
		log.Fatal("Database connection not initialized. Ensure db.Connect() is called in main.")
	}
	log.Println("Internal DB handle for boilerplate-service initialized successfully.")
}
