// File: services/website-cms-service/internal/repositories/item_repository_interface.go
// ---
package repositories

import "ape-go-services/website-cms-service/internal/models"

// ItemRepositoryInterface defines the contract for item repository operations.
// This allows for mocking in unit tests.
type ItemRepositoryInterface interface {
	GetAll() ([]models.Item, error)
	GetByID(id string) (*models.Item, error)
	Create(item *models.Item) error
	Update(id string, item *models.Item) error
	Delete(id string) error
	Count() (int64, error)
}
