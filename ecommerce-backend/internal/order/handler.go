// internal/order/handler.go
package order

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/utils" // For standardized responses
)

// OrderHandler handles HTTP requests related to orders.
type OrderHandler struct {
	Service   OrderService
	Validator *validator.Validate
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(s OrderService) *OrderHandler {
	return &OrderHandler{
		Service:   s,
		Validator: validator.New(),
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for the authenticated user
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param   request body CreateOrderRequest true "Order Creation Info"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{} "Order created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error, insufficient stock"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed: "+validationErrors.Error())
		return
	}

	// Get userID from Gin context, set by AuthMiddleware
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	userIDStr := userID.(string) // Cast to string

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second) // Longer timeout for transactions
	defer cancel()

	orderResp, err := h.Service.CreateOrder(ctx, userIDStr, &req)
	if err != nil {
		// Differentiate between user-facing errors (like insufficient stock) and internal errors
		if err.Error() == "invalid user ID format" ||
			(len(err.Error()) > 0 && (strings.Contains(err.Error(), "product not found") || strings.Contains(err.Error(), "insufficient stock"))) ||
			strings.Contains(err.Error(), "invalid product ID format") {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, gin.H{"message": "Order created successfully", "order": orderResp})
}

// GetUserOrders godoc
// @Summary Get orders for the authenticated user
// @Description Retrieve a list of all orders placed by the authenticated user
// @Tags Orders
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "List of user orders"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/my [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	userIDStr := userID.(string)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orders, err := h.Service.GetUserOrders(ctx, userIDStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"orders": orders})
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Retrieve a single order by its ID (accessible to user who placed it or admin)
// @Tags Orders
// @Produce  json
// @Param   id path string true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Order data"
// @Failure 400 {object} map[string]interface{} "Invalid order ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Access denied (if not owner or admin)"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.MustGet("userID").(string)     // User ID from authenticated context
	userRole := c.MustGet("userRole").(string) // User role from authenticated context

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderResp, err := h.Service.GetOrderByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "invalid order ID format" {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Authorization check: Only owner or admin can view this specific order
	if orderResp.UserID != userID && userRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Access denied: You can only view your own orders or if you are an admin.")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"order": orderResp})
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Retrieve a list of all orders (admin only)
// @Tags Orders
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "List of all orders"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (requires admin)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/orders [get]
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orders, err := h.Service.GetAllOrders(ctx)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"orders": orders})
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order by ID (admin only)
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param   id path string true "Order ID"
// @Param   request body UpdateOrderStatusRequest true "Order Status Update Info"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Order status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body, validation error, or invalid ID/status"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (requires admin)"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var req UpdateOrderStatusRequest
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

	orderResp, err := h.Service.UpdateOrderStatus(ctx, orderID, &req)
	if err != nil {
		if err.Error() == "order not found" {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "invalid order ID format" || strings.Contains(err.Error(), "invalid status") {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Order status updated successfully", "order": orderResp})
}
