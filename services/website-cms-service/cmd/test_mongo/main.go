package main

import (
	"ape-go-services/pkg/mongodb"
	"ape-go-services/pkg/repository"
	"ape-go-services/website-cms-service/internal/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil { // <-- REMOVE THE PATH
		log.Println("No .env file found, continuing with system env vars")
	}

	// Connect to MongoDB
	mongoConfig := mongodb.LoadConfigFromEnv()
	if err := mongodb.Connect(mongoConfig); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect()

	log.Println("=== MongoDB Connection Test ===")
	testConnection()

	log.Println("\n=== Repository Test ===")
	testRepository()

	log.Println("\n=== Service Layer Test ===")
	testService()

	log.Println("\n=== Test completed successfully! ===")
}

func testConnection() {
	if mongodb.IsConnected() {
		log.Println("✅ Successfully connected to MongoDB Atlas!")
		log.Printf("📊 Database: %s", os.Getenv("MONGO_DATABASE"))
	} else {
		log.Println("❌ Failed to connect to MongoDB")
		return
	}
}

func testRepository() {
	repo := repository.NewItemRepository()

	// Test retrieving all items
	log.Println("🔍 Retrieving all items from the database...")
	items, err := repo.GetAll()
	if err != nil {
		log.Printf("❌ Error retrieving items: %v", err)
		return
	}

	log.Printf("✅ Successfully retrieved %d items!", len(items))

	// Print items details
	for i, item := range items {
		log.Printf("📦 Item %d:", i+1)
		log.Printf("   ID: %s", item.ID.Hex())
		log.Printf("   Name: %s", item.Name)
		log.Printf("   Quantity: %d", item.Quantity)

		// Safely print CreatedAt
		if item.CreatedAt != nil {
			log.Printf("   Created: %s", item.CreatedAt.Format("2006-01-02 15:04:05"))
		} else {
			log.Printf("   Created: (not set)")
		}

		// Safely print UpdatedAt
		if item.UpdatedAt != nil {
			log.Printf("   Updated: %s", item.UpdatedAt.Format("2006-01-02 15:04:05"))
		} else {
			log.Printf("   Updated: (not set)")
		}

		log.Println()
	}

	// Test count
	count, err := repo.Count()
	if err != nil {
		log.Printf("❌ Error counting items: %v", err)
		return
	}
	log.Printf("📊 Total items in collection: %d", count)
}

func testService() {
	service := services.NewItemService()

	// Test service layer
	log.Println("🔄 Testing service layer...")
	items, err := service.GetAllItems()
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		return
	}

	log.Printf("✅ Service successfully retrieved %d items", len(items))

	// Test count through service
	count, err := service.GetItemCount()
	if err != nil {
		log.Printf("❌ Service count error: %v", err)
		return
	}
	log.Printf("📊 Service reports %d total items", count)
}
