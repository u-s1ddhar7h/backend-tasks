// internal/product/handler.go
package product

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"                  // For request body validation
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/utils" // For standardized responses
	"go.mongodb.org/mongo-driver/bson/primitive"              // For converting ID strings
)

// ProductHandler handles HTTP requests related to products.
type ProductHandler struct {
	Service   ProductService
	Validator *validator.Validate
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(s ProductService) *ProductHandler {
	return &ProductHandler{
		Service:   s,
		Validator: validator.New(),
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (admin only)
// @Tags Products
// @Accept  json
// @Produce  json
// @Param   request body ProductCreateRequest true "Product Create Info"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{} "Product created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 409 {object} map[string]interface{} "Product with this SKU already exists"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (requires admin)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed: "+validationErrors.Error())
		return
	}

	// Add context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productResp, err := h.Service.CreateProduct(ctx, &req)
	if err != nil {
		if err.Error() == "product with this SKU already exists" {
			utils.RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		if err.Error() == "invalid category ID format" {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, gin.H{"message": "Product created successfully", "product": productResp})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Retrieve a single product by its ID
// @Tags Products
// @Produce  json
// @Param   id path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Product data"
// @Failure 400 {object} map[string]interface{} "Invalid product ID format"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id") // Get ID from URL parameter

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productResp, err := h.Service.GetProductByID(ctx, productID)
	if err != nil {
		if err.Error() == "product not found" {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "invalid product ID format" {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"product": productResp})
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve a list of all products
// @Tags Products
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of products"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	products, err := h.Service.GetAllProducts(ctx)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"products": products})
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update an existing product by ID (admin only)
// @Tags Products
// @Accept  json
// @Produce  json
// @Param   id path string true "Product ID"
// @Param   request body ProductUpdateRequest true "Product Update Info"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Product updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body, validation error, or invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (requires admin)"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var req ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Manually validate CategoryID if it's present, as validator.v10 doesn't have an ObjectID rule directly
	if req.CategoryID != nil && *req.CategoryID != "" {
		if !primitive.IsValidObjectID(*req.CategoryID) {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid category ID format")
			return
		}
	}

	// Validate other fields that are pointers (omitempty validation)
	if err := h.Validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed: "+validationErrors.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productResp, err := h.Service.UpdateProduct(ctx, productID, &req)
	if err != nil {
		if err.Error() == "product not found" {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "invalid product ID format" || err.Error() == "no fields provided for update" {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Product updated successfully", "product": productResp})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID (admin only)
// @Tags Products
// @Produce  json
// @Param   id path string true "Product ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Product deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid product ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (requires admin)"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.Service.DeleteProduct(ctx, productID)
	if err != nil {
		if err.Error() == "product not found" {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "invalid product ID format" {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
