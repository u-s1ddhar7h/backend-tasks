// cmd/api/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/auth"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/config"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/database"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/middleware"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/order"   // NEW: Import order package
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/product" // Import product package
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Connect to Database
	database.ConnectDB(cfg.MongoURI)
	defer func() {
		if database.MongoClient != nil {
			log.Println("Disconnecting from MongoDB...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := database.MongoClient.Disconnect(ctx); err != nil {
				log.Fatalf("Error disconnecting from MongoDB: %v", err)
			}
			log.Println("MongoDB disconnected.")
		}
	}()

	// 3. Initialize Gin Router
	router := gin.Default()

	// 4. Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Be specific in production, e.g., "http://localhost:3000"
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5. Initialize Services and Handlers
	authService := auth.NewAuthService(cfg)
	authHandler := auth.NewAuthHandler(authService)

	productService := product.NewProductService()
	productHandler := product.NewProductHandler(productService)

	// NEW: Initialize Order Service and Handler.
	// OrderService needs ProductService injected because it interacts with product stock.
	orderService := order.NewOrderService(productService)
	orderHandler := order.NewOrderHandler(orderService)

	// 6. Define Routes
	// Public routes (no authentication required)
	publicRoutes := router.Group("/api")
	{
		// Authentication routes
		publicRoutes.POST("/auth/register", authHandler.Register)
		publicRoutes.POST("/auth/login", authHandler.Login)

		// Public product routes (view products without login)
		publicRoutes.GET("/products", productHandler.GetAllProducts)
		publicRoutes.GET("/products/:id", productHandler.GetProductByID)
	}

	// Authenticated routes (require a valid JWT)
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware(cfg)) // Apply the authentication middleware
	{
		protectedRoutes.GET("/auth/me", authHandler.GetMe)

		// Admin-only product routes (create, update, delete)
		adminProducts := protectedRoutes.Group("/products")
		adminProducts.Use(middleware.AuthorizeRole("admin")) // Requires "admin" role
		{
			adminProducts.POST("/", productHandler.CreateProduct)
			adminProducts.PUT("/:id", productHandler.UpdateProduct)
			adminProducts.DELETE("/:id", productHandler.DeleteProduct)
		}

		// User-authenticated order routes
		userOrders := protectedRoutes.Group("/orders")
		{
			userOrders.POST("/", orderHandler.CreateOrder)    // Create a new order
			userOrders.GET("/my", orderHandler.GetUserOrders) // Get all orders for the authenticated user
			userOrders.GET("/:id", orderHandler.GetOrderByID) // Get a specific order (with ownership/admin check inside handler)
		}

		// Admin-only order routes
		adminOrders := protectedRoutes.Group("/admin/orders")
		adminOrders.Use(middleware.AuthorizeRole("admin")) // Requires "admin" role
		{
			adminOrders.GET("/", orderHandler.GetAllOrders)                  // Get all orders in the system
			adminOrders.PATCH("/:id/status", orderHandler.UpdateOrderStatus) // Update order status
		}
	}

	// 7. Start the server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
