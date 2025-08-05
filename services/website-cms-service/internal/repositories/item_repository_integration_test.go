// File: services/website-cms-service/internal/repositories/item_repository_integration_test.go
// ---
//go:build integration

package repositories

import (
	"ape-go-services/pkg/db"
	"ape-go-services/website-cms-service/internal/models"
	"context"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ItemRepositoryIntegrationSuite is a test suite for the ItemRepository.
// It uses the testify suite package to manage setup and teardown of the database connection.
type ItemRepositoryIntegrationSuite struct {
	suite.Suite
	repo    *ItemRepository
	timeout time.Duration
}

// SetupSuite runs once before any tests in the suite are run.
// It's responsible for establishing the database connection.
func (s *ItemRepositoryIntegrationSuite) SetupSuite() {
	// The test is in a sub-package, so we need to go up to the service root to find the .env file.
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file for integration tests: %v", err)
	}

	mongoConfig := db.LoadMongoConfigFromEnv()
	err := db.ConnectMongoDB(mongoConfig)
	s.Require().NoError(err, "Failed to connect to MongoDB for integration tests")

	s.repo = NewItemRepository()
	s.timeout = 30 * time.Second
	log.Println("MongoDB connection established for integration test suite.")
}

// TearDownSuite runs once after all tests in the suite have finished.
// It's responsible for closing the database connection.
func (s *ItemRepositoryIntegrationSuite) TearDownSuite() {
	if err := db.DisconnectMongoDB(); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
	log.Println("MongoDB connection closed.")
}

// SetupTest runs before each individual test function.
// It cleans the collection to ensure tests run in isolation and don't affect each other.
func (s *ItemRepositoryIntegrationSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err := s.repo.collection.DeleteMany(ctx, bson.M{})
	s.Require().NoError(err, "Failed to clean the items collection before a test")
}

// TestItemRepositoryIntegrationSuite is the entry point for the test runner.
func TestItemRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ItemRepositoryIntegrationSuite))
}

// --- Integration Test Cases ---

func (s *ItemRepositoryIntegrationSuite) TestCreateAndGetByID() {
	t := s.T()

	// Arrange: Create a new item model.
	newItem := models.NewItem("MacBook Pro", 5)

	// Act: Call the Create method.
	err := s.repo.Create(newItem)
	s.Require().NoError(err, "Create should not return an error")
	s.Require().NotEqual(primitive.NilObjectID, newItem.ID, "ID should be set by the database after creation")

	// Assert: Retrieve the item from the DB and check its fields.
	retrievedItem, err := s.repo.GetByID(newItem.ID.Hex())
	assert.NoError(t, err, "GetByID should not return an error for a valid ID")
	assert.NotNil(t, retrievedItem, "Retrieved item should not be nil")
	assert.Equal(t, "MacBook Pro", retrievedItem.Name)
	assert.Equal(t, 5, retrievedItem.Quantity)
	assert.WithinDuration(t, time.Now(), *retrievedItem.CreatedAt, 10*time.Second)
}

func (s *ItemRepositoryIntegrationSuite) TestGetAll() {
	t := s.T()

	// Arrange: Create multiple items.
	item1 := models.NewItem("Wireless Mouse", 25)
	item2 := models.NewItem("Mechanical Keyboard", 15)
	s.Require().NoError(s.repo.Create(item1))
	s.Require().NoError(s.repo.Create(item2))

	// Act: Get all items.
	items, err := s.repo.GetAll()

	// Assert: Check that the correct number of items were returned.
	assert.NoError(t, err)
	assert.Len(t, items, 2, "GetAll should return all created items")
}

func (s *ItemRepositoryIntegrationSuite) TestUpdate() {
	t := s.T()

	// Arrange: Create an initial item.
	item := models.NewItem("Old Product Name", 10)
	s.Require().NoError(s.repo.Create(item))

	// Act: Update the item's fields.
	item.Name = "New & Improved Product"
	item.Quantity = 100
	err := s.repo.Update(item.ID.Hex(), item)

	// Assert: Check that the update was successful and the fields are changed in the DB.
	assert.NoError(t, err, "Update should succeed")
	updatedItem, _ := s.repo.GetByID(item.ID.Hex())
	assert.Equal(t, "New & Improved Product", updatedItem.Name)
	assert.Equal(t, 100, updatedItem.Quantity)
	assert.True(t, updatedItem.UpdatedAt.After(*updatedItem.CreatedAt), "UpdatedAt should be more recent than CreatedAt")
}

func (s *ItemRepositoryIntegrationSuite) TestDelete() {
	t := s.T()

	// Arrange: Create an item to be deleted.
	item := models.NewItem("Temporary Item", 1)
	s.Require().NoError(s.repo.Create(item))
	id := item.ID.Hex()

	// Act: Delete the item.
	err := s.repo.Delete(id)

	// Assert: Ensure it was deleted and can no longer be found.
	assert.NoError(t, err, "Delete should succeed")
	_, err = s.repo.GetByID(id)
	assert.Error(t, err, "GetByID should return an error for a deleted item")
	assert.Contains(t, err.Error(), "item not found")
}

func (s *ItemRepositoryIntegrationSuite) TestCount() {
	t := s.T()

	// Arrange: Create a known number of items.
	s.repo.Create(models.NewItem("A", 1))
	s.repo.Create(models.NewItem("B", 1))
	s.repo.Create(models.NewItem("C", 1))

	// Act: Get the count.
	count, err := s.repo.Count()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count, "Count should return the correct number of items")
}
