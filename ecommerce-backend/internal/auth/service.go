// internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/u-s1ddhar7h/ecommerce-backend/internal/config" // Import config to get JWT_SECRET
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/database"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/utils" // For password hashing and JWT
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	RegisterUser(ctx context.Context, req *RegisterRequest) (string, *UserResponse, error)
	LoginUser(ctx context.Context, req *LoginRequest) (string, *UserResponse, error)
	GetUserByID(ctx context.Context, userID string) (*UserResponse, error)
}

// service implements AuthService.
type service struct {
	usersCollection *mongo.Collection
	cfg             *config.Config // Store config to access JWTSecret
}

// NewAuthService creates a new authentication service.
func NewAuthService(cfg *config.Config) AuthService {
	return &service{
		usersCollection: database.GetCollection("users"), // Get the 'users' collection
		cfg:             cfg,
	}
}

// RegisterUser handles new user registration.
func (s *service) RegisterUser(ctx context.Context, req *RegisterRequest) (string, *UserResponse, error) {
	// Check if user with this email already exists
	var existingUser User
	err := s.usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return "", nil, errors.New("user with this email already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking for existing user: %v", err)
		return "", nil, errors.New("database error during registration check")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", nil, errors.New("failed to hash password")
	}

	// Create new user object
	now := time.Now()
	user := &User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "user", // Default role
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert user into MongoDB
	result, err := s.usersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		return "", nil, errors.New("failed to register user")
	}

	user.ID = result.InsertedID.(primitive.ObjectID) // Set the generated ID

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Role, s.cfg.JWTSecret)
	if err != nil {
		log.Printf("Error generating JWT for new user: %v", err)
		return "", nil, errors.New("failed to generate authentication token")
	}

	// Prepare response data (without password)
	userResp := &UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return token, userResp, nil
}

// LoginUser handles user login.
func (s *service) LoginUser(ctx context.Context, req *LoginRequest) (string, *UserResponse, error) {
	var user User
	err := s.usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil, errors.New("invalid credentials")
		}
		log.Printf("Error finding user during login: %v", err)
		return "", nil, errors.New("database error during login")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Role, s.cfg.JWTSecret)
	if err != nil {
		log.Printf("Error generating JWT for login: %v", err)
		return "", nil, errors.New("failed to generate authentication token")
	}

	// Prepare response data (without password)
	userResp := &UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return token, userResp, nil
}

// GetUserByID retrieves a user by their ID for authenticated endpoints.
func (s *service) GetUserByID(ctx context.Context, userID string) (*UserResponse, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user User
	err = s.usersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		log.Printf("Error retrieving user by ID: %v", err)
		return nil, errors.New("database error retrieving user")
	}

	userResp := &UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResp, nil
}
