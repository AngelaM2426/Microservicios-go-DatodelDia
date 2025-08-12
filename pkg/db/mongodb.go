// File: pkg/db/mongodb.go
package db

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

// Renamed global variables to be specific to MongoDB
var (
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
)

// MongoConfig holds MongoDB configuration. The name is unchanged as it's a local struct.
type MongoConfig struct {
	URI         string
	Database    string
	Timeout     time.Duration
	MaxPoolSize uint64
}

// LoadMongoConfigFromEnv loads MongoDB configuration from environment variables.
// Renamed from LoadConfigFromEnv -> LoadMongoConfigFromEnv
func LoadMongoConfigFromEnv() *MongoConfig {
	return &MongoConfig{
		URI:         getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		Database:    getEnvOrDefault("MONGO_DATABASE", "myFirstDatabase"),
		Timeout:     30 * time.Second,
		MaxPoolSize: 100,
	}
}

// ConnectMongoDB establishes a MongoDB connection using the provided config.
// Renamed from Connect -> ConnectMongoDB
func ConnectMongoDB(config *MongoConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Set global variables (which have also been renamed)
	MongoClient = client
	MongoDatabase = client.Database(config.Database)

	log.Printf("Successfully connected to MongoDB database: %s", config.Database)
	return nil
}

// DisconnectMongoDB closes the MongoDB connection.
// Renamed from Disconnect -> DisconnectMongoDB
func DisconnectMongoDB() error {
	if MongoClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := MongoClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

// GetMongoCollection returns a MongoDB collection.
// Renamed from GetCollection -> GetMongoCollection
func GetMongoCollection(name string) *mongo.Collection {
	if MongoDatabase == nil {
		log.Fatal("MongoDB Database not initialized. Call ConnectMongoDB() first.")
	}
	return MongoDatabase.Collection(name)
}

// IsMongoDBConnected checks if the MongoDB connection is active.
// Renamed from IsConnected -> IsMongoDBConnected
func IsMongoDBConnected() bool {
	if MongoClient == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return MongoClient.Ping(ctx, readpref.Primary()) == nil
}

// getEnvOrDefault is a private helper function, no changes needed here.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
