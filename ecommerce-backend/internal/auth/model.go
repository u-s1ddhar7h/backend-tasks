// internal/auth/model.go
package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system.
// bson tags are for MongoDB document mapping.
// json tags are for HTTP request/response JSON serialization.
// validate tags are for request body validation (from github.com/go-playground/validator).
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB's unique ID for the document
	Username  string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"` // Hashed password
	Role      string             `bson:"role" json:"role"`                                   // e.g., "user", "admin"
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// LoginRequest defines the structure for a login request body.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest defines the structure for a registration request body.
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse defines the structure for a user's data sent in API responses (e.g., /me endpoint)
// It deliberately omits the password field for security.
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
