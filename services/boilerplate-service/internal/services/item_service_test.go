// File: services/boilerplate-service/internal/services/item_service_test.go
// ---
package services

import (
	"ape-go-services/boilerplate-service/internal/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockItemRepository is a mock implementation of ItemRepositoryInterface
type MockItemRepository struct {
	mock.Mock
}

// Implement the interface methods for the mock
func (m *MockItemRepository) GetAll() ([]models.Item, error) {
	args := m.Called()
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *MockItemRepository) GetByID(id string) (*models.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

func (m *MockItemRepository) Create(item *models.Item) error {
	args := m.Called(item)
	item.ID = primitive.NewObjectID()
	return args.Error(0)
}

func (m *MockItemRepository) Update(id string, item *models.Item) error {
	args := m.Called(id, item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockItemRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// --- Unit Tests for ItemService ---

func TestItemService_GetAllItems_Success(t *testing.T) {
	// Arrange
	t.Logf("Arrange: Setting up mock repository and expected items.")
	mockRepo := new(MockItemRepository)
	itemService := NewItemServiceWithRepository(mockRepo)
	now := time.Now()
	expectedItems := []models.Item{
		{ID: primitive.NewObjectID(), Name: "Item 1", Quantity: 10, CreatedAt: &now},
		{ID: primitive.NewObjectID(), Name: "Item 2", Quantity: 20, CreatedAt: &now},
	}
	mockRepo.On("GetAll").Return(expectedItems, nil)

	// Act
	t.Logf("Act: Calling GetAllItems() service method.")
	items, err := itemService.GetAllItems()

	// Assert
	t.Logf("Assert: Verifying the results. Expecting %d items.", len(expectedItems))
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	assert.Len(t, items, 2)
	mockRepo.AssertExpectations(t)
}

func TestItemService_GetAllItems_RepositoryError(t *testing.T) {
	// Arrange
	t.Logf("Arrange: Setting up mock repository to return an error.")
	mockRepo := new(MockItemRepository)
	itemService := NewItemServiceWithRepository(mockRepo)
	expectedError := errors.New("database is down")
	mockRepo.On("GetAll").Return([]models.Item{}, expectedError)

	// Act
	t.Logf("Act: Calling GetAllItems(), expecting an error.")
	items, err := itemService.GetAllItems()

	// Assert
	t.Logf("Assert: Verifying that an error was returned.")
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Contains(t, err.Error(), "service: failed to get items")
	mockRepo.AssertExpectations(t)
}

// TestItemService_GetByID uses t.Run to create structured sub-tests
func TestItemService_GetByID(t *testing.T) {
	// Arrange: This setup is shared for all sub-tests
	mockRepo := new(MockItemRepository)
	itemService := NewItemServiceWithRepository(mockRepo)
	now := time.Now()
	testID := primitive.NewObjectID()
	expectedItem := &models.Item{
		ID: testID, Name: "Found Item", Quantity: 99, CreatedAt: &now,
	}

	// Setup mock expectations
	mockRepo.On("GetByID", testID.Hex()).Return(expectedItem, nil)
	mockRepo.On("GetByID", "nonexistent_id").Return(nil, errors.New("item not found"))

	// Sub-test for the success case
	t.Run("Success case - item found", func(t *testing.T) {
		// Act
		t.Logf("Act: Calling GetItemByID with a valid ID: %s", testID.Hex())
		item, err := itemService.GetItemByID(testID.Hex()) // <--- CORRECTED HERE

		// Assert
		t.Logf("Assert: Verifying the correct item is returned.")
		assert.NoError(t, err)
		assert.Equal(t, expectedItem, item)
	})

	// Sub-test for the failure case
	t.Run("Failure case - item not found", func(t *testing.T) {
		// Act
		t.Logf("Act: Calling GetItemByID with an invalid ID.")
		item, err := itemService.GetItemByID("nonexistent_id") // <--- CORRECTED HERE

		// Assert
		t.Logf("Assert: Verifying a 'not found' error is returned.")
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "item not found")
	})

	// Verify that all expected mock calls were made
	mockRepo.AssertExpectations(t)
}

func TestItemService_CreateItem_Success(t *testing.T) {
	// Arrange
	t.Logf("Arrange: Setting up mock to accept a new item.")
	mockRepo := new(MockItemRepository)
	itemService := NewItemServiceWithRepository(mockRepo)
	mockRepo.On("Create", mock.AnythingOfType("*models.Item")).Return(nil)

	// Act
	t.Logf("Act: Calling CreateItem with name 'New Gadget' and quantity 15.")
	createdItem, err := itemService.CreateItem("New Gadget", 15)

	// Assert
	t.Logf("Assert: Verifying a new item was returned with a generated ID.")
	assert.NoError(t, err)
	assert.NotNil(t, createdItem)
	assert.Equal(t, "New Gadget", createdItem.Name)
	assert.NotEqual(t, primitive.NilObjectID, createdItem.ID)
	mockRepo.AssertExpectations(t)
}

func TestItemService_GetItemCount_Success(t *testing.T) {
	// Arrange
	t.Logf("Arrange: Mocking the Count() method to return 42.")
	mockRepo := new(MockItemRepository)
	itemService := NewItemServiceWithRepository(mockRepo)
	mockRepo.On("Count").Return(int64(42), nil)

	// Act
	t.Logf("Act: Calling GetItemCount().")
	count, err := itemService.GetItemCount()

	// Assert
	t.Logf("Assert: Verifying count is 42.")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), count)
	mockRepo.AssertExpectations(t)
}
