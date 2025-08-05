package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM database instance
var DB *gorm.DB

// sqlDB is the underlying sql.DB instance, kept private to this package
var sqlDB *gorm.DB

// Connect establishes a PostgreSQL connection, tests it, and prepares it for use.
func Connect() (*gorm.DB, error) {
	// --- Configuration Loading ---
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Use logger.Info for verbose query logging
	}

	// --- Open Connection ---
	// gorm.Open does not establish a connection, it only prepares the config.
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// --- Connection Pool Configuration ---
	// Get the underlying sql.DB object to configure the connection pool.
	rawDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool parameters as needed.
	// TODO: Adjust these values based on your application's requirements.
	// TODO: Consider making these configurable via environment variables.
	rawDB.SetMaxOpenConns(25)
	rawDB.SetMaxIdleConns(10)
	rawDB.SetConnMaxLifetime(time.Hour)

	// --- Connection Test (Ping) ---
	// Ping the database to ensure the connection is actually alive.
	log.Println("Pinging PostgreSQL database...")
	if err := rawDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	log.Printf("Successfully connected to Postgres database: %s", dbname)

	// --- Set Global Variables ---
	// Store the connection instances for global access and for the Disconnect function.
	DB = db
	sqlDB = db // Store the gorm.DB instance to access its .DB() method later.

	return DB, nil
}

// Disconnect gracefully closes the PostgreSQL database connection.
func Disconnect() error {
	if sqlDB == nil {
		log.Println("No PostgreSQL connection to close.")
		return nil
	}

	log.Println("Closing PostgreSQL connection...")
	rawDB, err := sqlDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for closing: %w", err)
	}

	if err := rawDB.Close(); err != nil {
		return fmt.Errorf("failed to close postgres connection: %w", err)
	}

	log.Println("PostgreSQL connection closed gracefully.")
	return nil
}

// GetDB returns the active GORM database instance.
// This is a helper to ensure we don't need to import and manage the DB instance themselves.
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("PostgreSQL connection has not been initialized. Call db.Connect() first.")
	}
	return DB
}
