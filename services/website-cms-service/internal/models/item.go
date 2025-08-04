// File: services/website-cms-service/internal/models/item.go
// ---
package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Item represents an item document in MongoDB
type Item struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	CreatedAt *time.Time         `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// NewItem creates a new Item with timestamps
func NewItem(name string, quantity int) *Item {
	now := time.Now()
	return &Item{
		Name:      name,
		Quantity:  quantity,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

// ToMap converts Item to a map for updates
func (i *Item) ToMap() bson.M {
	return bson.M{
		"name":       i.Name,
		"quantity":   i.Quantity,
		"updated_at": time.Now(),
	}
}

// Validate validates the Item fields
func (i *Item) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("name is required")
	}
	if i.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}
	return nil
}
