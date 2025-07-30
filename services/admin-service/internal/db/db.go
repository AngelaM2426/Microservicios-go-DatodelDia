package db

import (
	"ape-go-services/pkg/db"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	// Connect to the User Service database
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = conn
	log.Println("Database connected successfully")
}
