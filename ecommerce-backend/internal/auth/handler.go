// internal/auth/handler.go
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/utils" // For response helpers
)

// AuthHandler handles HTTP requests related to authentication.
type AuthHandler struct {
	Service   AuthService
	Validator *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{
		Service:   s,
		Validator: validator.New(), // Initialize validator
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body RegisterRequest true "Registration Info"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 409 {object} map[string]interface{} "User with this email already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	// Bind JSON request body to the RegisterRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate the request struct fields
	if err := h.Validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed: "+validationErrors.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Use Gin's request context
	defer cancel()

	token, userResp, err := h.Service.RegisterUser(ctx, &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, gin.H{"message": "User registered successfully", "token": token, "user": userResp})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate a user with email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body LoginRequest true "Login Info"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed: "+validationErrors.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	token, userResp, err := h.Service.LoginUser(ctx, &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			utils.RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Login successful", "token": token, "user": userResp})
}

// GetMe godoc
// @Summary Get current user's profile
// @Description Retrieve the profile of the authenticated user
// @Tags Auth
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "User profile data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Retrieve userID from Gin context, which was set by AuthMiddleware
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "User ID not found in context (middleware error)")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert userID to string if it's not already (it should be string from JWT claims)
	userIDStr, ok := userID.(string)
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "Invalid userID type in context")
		return
	}

	userResp, err := h.Service.GetUserByID(ctx, userIDStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"user": userResp})
}
