// internal/order/model.go
package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItem represents a single product within an order.
type OrderItem struct {
	ProductID primitive.ObjectID `bson:"productID" json:"productId"`
	Name      string             `bson:"name" json:"name"` // Denormalized product name
	SKU       string             `bson:"sku" json:"sku"`   // Denormalized product SKU
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"` // Price at time of order
	Subtotal  float64            `bson:"subtotal" json:"subtotal"`
}

// Order represents a customer order.
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID" json:"userId"`
	Items       []OrderItem        `bson:"items" json:"items"`
	TotalAmount float64            `bson:"totalAmount" json:"totalAmount"`
	Status      string             `bson:"status" json:"status"` // e.g., "pending", "processing", "shipped", "delivered", "cancelled"
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// OrderStatus defines possible statuses for an order.
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusShipped    = "shipped"
	StatusDelivered  = "delivered"
	StatusCancelled  = "cancelled"
)

// CreateOrderRequest defines the structure for a new order request body.
type CreateOrderRequest struct {
	Items []struct {
		ProductID string `json:"productId" validate:"required"`
		Quantity  int    `json:"quantity" validate:"required,gt=0"`
	} `json:"items" validate:"required,min=1,dive"` // `dive` validates each item in the slice
}

// UpdateOrderStatusRequest defines the structure for updating an order's status.
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending processing shipped delivered cancelled"`
}

// OrderResponse defines the structure for order data in API responses.
type OrderResponse struct {
	ID          string      `json:"id"`
	UserID      string      `json:"userId"`
	Items       []OrderItem `json:"items"` // Items are typically fine to return as is
	TotalAmount float64     `json:"totalAmount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
