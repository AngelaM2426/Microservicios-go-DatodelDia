// File: services/website-cms-service/internal/services/item_service.go
// ---
package services

import (
	"ape-go-services/website-cms-service/internal/models"
	"ape-go-services/website-cms-service/internal/repositories"
	"fmt"
	"log"
)

// ItemService handles business logic for items
type ItemService struct {
	repo repositories.ItemRepositoryInterface
}

// NewItemService creates a new ItemService with a real repository.
func NewItemService() *ItemService {
	return &ItemService{
		repo: repositories.NewItemRepository(),
	}
}

// NewItemServiceWithRepository creates a new ItemService with a custom repository for testing.
func NewItemServiceWithRepository(repo repositories.ItemRepositoryInterface) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

// GetAllItems retrieves all items with business logic
func (s *ItemService) GetAllItems() ([]models.Item, error) {
	items, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("service: failed to get items: %w", err)
	}

	log.Printf("Retrieved %d items from database", len(items))
	return items, nil
}

// GetItemByID retrieves a single item by ID  <-- THIS IS THE NEW METHOD
func (s *ItemService) GetItemByID(id string) (*models.Item, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get item: %w", err)
	}

	return item, nil
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(name string, quantity int) (*models.Item, error) {
	item := models.NewItem(name, quantity)

	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("service: failed to create item: %w", err)
	}

	log.Printf("Created new item: %s (ID: %s)", item.Name, item.ID.Hex())
	return item, nil
}

// UpdateItem updates an existing item
func (s *ItemService) UpdateItem(id string, name string, quantity int) error {
	item := &models.Item{
		Name:     name,
		Quantity: quantity,
	}

	if err := s.repo.Update(id, item); err != nil {
		return fmt.Errorf("service: failed to update item: %w", err)
	}

	log.Printf("Updated item with ID: %s", id)
	return nil
}

// DeleteItem deletes an item
func (s *ItemService) DeleteItem(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("service: failed to delete item: %w", err)
	}

	log.Printf("Deleted item with ID: %s", id)
	return nil
}

// GetItemCount returns the total number of items
func (s *ItemService) GetItemCount() (int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return 0, fmt.Errorf("service: failed to count items: %w", err)
	}

	return count, nil
}
