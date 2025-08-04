package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

// Config holds MongoDB configuration
type Config struct {
	URI         string
	Database    string
	Timeout     time.Duration
	MaxPoolSize uint64
}

// LoadConfigFromEnv loads MongoDB configuration from environment variables
func LoadConfigFromEnv() *Config {
	return &Config{
		URI:         getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		Database:    getEnvOrDefault("MONGO_DATABASE", "myFirstDatabase"),
		Timeout:     30 * time.Second,
		MaxPoolSize: 100,
	}
}

// Connect establishes a MongoDB connection using the provided config
func Connect(config *Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Client options
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(config.Timeout)

	// Connect to MongoDB
	// client, err := mongo.Connect(clientOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Set global variables
	Client = client
	Database = client.Database(config.Database)

	log.Printf("Successfully connected to MongoDB database: %s", config.Database)
	return nil
}

// Disconnect closes the MongoDB connection
func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(name string) *mongo.Collection {
	if Database == nil {
		log.Fatal("Database not initialized. Call Connect() first.")
	}
	return Database.Collection(name)
}

// IsConnected checks if the MongoDB connection is active
func IsConnected() bool {
	if Client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return Client.Ping(ctx, readpref.Primary()) == nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
