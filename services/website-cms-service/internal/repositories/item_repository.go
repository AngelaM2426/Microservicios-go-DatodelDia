// File: services/website-cms-service/internal/repositories/item_repository.go
// ---
package repositories

import (
	"context"
	"fmt"
	"time"

	"ape-go-services/pkg/mongodb"
	"ape-go-services/website-cms-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ItemsCollectionName = "items"

// ItemRepository implements the data access logic for items.
type ItemRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewItemRepository creates a new ItemRepository
func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		collection: mongodb.GetCollection(ItemsCollectionName),
		timeout:    30 * time.Second,
	}
}

// GetAll retrieves all items from the database
func (r *ItemRepository) GetAll() ([]models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find items: %w", err)
	}
	defer cursor.Close(ctx)

	var items []models.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode items: %w", err)
	}

	return items, nil
}

// GetByID retrieves an item by its ID
func (r *ItemRepository) GetByID(id string) (*models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var item models.Item
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	return &item, nil
}

// Create inserts a new item into the database
func (r *ItemRepository) Create(item *models.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if err := item.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	now := time.Now()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update modifies an existing item in the database
func (r *ItemRepository) Update(id string, item *models.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	if err := item.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	update := bson.M{
		"$set": item.ToMap(),
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// Delete removes an item from the database
func (r *ItemRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// Count returns the total number of items
func (r *ItemRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count items: %w", err)
	}

	return count, nil
}
