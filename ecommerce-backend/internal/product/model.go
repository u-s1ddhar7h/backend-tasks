// internal/product/model.go
package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the system.
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required,min=3,max=100"`
	Description string             `bson:"description" json:"description" validate:"required,min=10,max=500"`
	Price       float64            `bson:"price" json:"price" validate:"required,gt=0"`              // gt=0 means greater than 0
	SKU         string             `bson:"sku" json:"sku" validate:"required,alphanum,min=5,max=20"` // Stock Keeping Unit
	CategoryID  primitive.ObjectID `bson:"categoryID" json:"categoryID" validate:"required"`         // Reference to the Category
	Stock       int                `bson:"stock" json:"stock" validate:"required,gte=0"`             // gte=0 means greater than or equal to 0
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// ProductResponse defines the structure for product data in API responses.
// It might be slightly different from the internal Product struct.
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
	CategoryID  string    `json:"categoryID"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ProductCreateRequest defines the structure for creating a new product.
type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	SKU         string  `json:"sku" validate:"required,alphanum,min=5,max=20"`
	CategoryID  string  `json:"categoryID" validate:"required"` // We expect the CategoryID as a string from the request
	Stock       int     `json:"stock" validate:"required,gte=0"`
}

// ProductUpdateRequest defines the structure for updating an existing product.
// All fields are optional, so we can update only specific fields.
type ProductUpdateRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=100"` // Pointers to allow optional fields
	Description *string  `json:"description,omitempty" validate:"omitempty,min=10,max=500"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	SKU         *string  `json:"sku,omitempty" validate:"omitempty,alphanum,min=5,max=20"`
	CategoryID  *string  `json:"categoryID,omitempty"` // Optional
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
}
