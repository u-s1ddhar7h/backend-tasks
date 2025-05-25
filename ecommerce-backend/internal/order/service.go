// internal/order/service.go
package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/u-s1ddhar7h/ecommerce-backend/internal/database"
	"github.com/u-s1ddhar7h/ecommerce-backend/internal/product" // Import product service to get product details and update stock
)

// OrderService defines the interface for order operations.
type OrderService interface {
	CreateOrder(ctx context.Context, userID string, req *CreateOrderRequest) (*OrderResponse, error)
	GetUserOrders(ctx context.Context, userID string) ([]OrderResponse, error)
	GetOrderByID(ctx context.Context, orderID string) (*OrderResponse, error)
	GetAllOrders(ctx context.Context) ([]OrderResponse, error)                                                    // Admin only
	UpdateOrderStatus(ctx context.Context, orderID string, req *UpdateOrderStatusRequest) (*OrderResponse, error) // Admin only
}

// service implements OrderService.
type service struct {
	ordersCollection *mongo.Collection
	productService   product.ProductService // Dependency on ProductService
}

// NewOrderService creates a new order service.
func NewOrderService(prodService product.ProductService) OrderService {
	return &service{
		ordersCollection: database.GetCollection("orders"), // Get the 'orders' collection
		productService:   prodService,
	}
}

// orderToResponse converts an Order model to an OrderResponse.
func orderToResponse(o *Order) *OrderResponse {
	return &OrderResponse{
		ID:          o.ID.Hex(),
		UserID:      o.UserID.Hex(),
		Items:       o.Items, // OrderItem already has json tags
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

// CreateOrder handles the creation of a new order.
func (s *service) CreateOrder(ctx context.Context, userID string, req *CreateOrderRequest) (*OrderResponse, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var orderItems []OrderItem
	var totalAmount float64

	// Declare 'order' outside the transaction closure
	var order Order
	var insertedID primitive.ObjectID // To hold the ID inserted by MongoDB

	// Start a MongoDB session for transaction
	session, err := database.MongoClient.StartSession()
	if err != nil {
		log.Printf("Error starting MongoDB session: %v", err)
		return nil, errors.New("failed to start database session")
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// Start transaction
		if err = session.StartTransaction(); err != nil {
			return err
		}

		for _, itemReq := range req.Items {
			productObjID, err := primitive.ObjectIDFromHex(itemReq.ProductID)
			if err != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("invalid product ID format for item %s", itemReq.ProductID)
			}

			productData, err := s.productService.GetProductForOrder(sessionContext, productObjID.Hex())
			if err != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("product not found or error retrieving product %s: %v", itemReq.ProductID, err)
			}

			if productData.Stock < itemReq.Quantity {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("insufficient stock for product '%s'. Available: %d, Requested: %d",
					productData.Name, productData.Stock, itemReq.Quantity)
			}

			updateResult, err := database.GetCollection("products").UpdateOne(
				sessionContext,
				bson.M{"_id": productObjID, "stock": bson.M{"$gte": itemReq.Quantity}},
				bson.M{"$inc": bson.M{"stock": -itemReq.Quantity}, "$set": bson.M{"updatedAt": time.Now()}},
			)
			if err != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("failed to deduct stock for product %s: %v", productData.Name, err)
			}
			if updateResult.MatchedCount == 0 || updateResult.ModifiedCount == 0 {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("failed to deduct stock for product %s (concurrent modification or insufficient stock after check)", productData.Name)
			}

			itemSubtotal := productData.Price * float64(itemReq.Quantity)
			orderItems = append(orderItems, OrderItem{
				ProductID: productData.ID,
				Name:      productData.Name,
				SKU:       productData.SKU,
				Quantity:  itemReq.Quantity,
				Price:     productData.Price,
				Subtotal:  itemSubtotal,
			})
			totalAmount += itemSubtotal
		}

		now := time.Now()
		// Assign to the 'order' variable declared outside
		order = Order{ // Note: assignment using '=' not ':=', and it's a struct, not a pointer initially
			UserID:      userObjectID,
			Items:       orderItems,
			TotalAmount: totalAmount,
			Status:      StatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		result, err := s.ordersCollection.InsertOne(sessionContext, &order) // Pass pointer for insertion
		if err != nil {
			session.AbortTransaction(sessionContext)
			log.Printf("Error inserting new order: %v", err)
			return errors.New("failed to create order")
		}

		insertedID = result.InsertedID.(primitive.ObjectID) // Store the inserted ID

		if err = session.CommitTransaction(sessionContext); err != nil {
			log.Printf("Error committing transaction: %v", err)
			return errors.New("failed to finalize order (transaction commit failed)")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Now 'insertedID' is accessible here
	// We can use GetOrderByID with this ID, or just construct the response from the 'order' variable
	// if we trust its state after the transaction.
	// For simplicity and to avoid another DB call if the transaction was complex:
	// Set the ID back to the order object.
	order.ID = insertedID

	return orderToResponse(&order), nil // Return the 'order' converted to response
}

// GetUserOrders retrieves all orders for a specific user.
func (s *service) GetUserOrders(ctx context.Context, userID string) ([]OrderResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	cursor, err := s.ordersCollection.Find(ctx, bson.M{"userID": userObjID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		log.Printf("Error finding user orders: %v", err)
		return nil, errors.New("failed to retrieve user orders")
	}
	defer cursor.Close(ctx)

	var orders []Order
	if err = cursor.All(ctx, &orders); err != nil {
		log.Printf("Error decoding user orders from cursor: %v", err)
		return nil, errors.New("failed to process user order data")
	}

	var orderResponses []OrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, *orderToResponse(&o))
	}

	return orderResponses, nil
}

// GetOrderByID retrieves a single order by its ID.
func (s *service) GetOrderByID(ctx context.Context, orderID string) (*OrderResponse, error) {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID format")
	}

	var order Order
	err = s.ordersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		log.Printf("Error finding order by ID: %v", err)
		return nil, errors.New("database error retrieving order")
	}

	return orderToResponse(&order), nil
}

// GetAllOrders retrieves all orders (for admin dashboard).
func (s *service) GetAllOrders(ctx context.Context) ([]OrderResponse, error) {
	cursor, err := s.ordersCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		log.Printf("Error finding all orders: %v", err)
		return nil, errors.New("failed to retrieve all orders")
	}
	defer cursor.Close(ctx)

	var orders []Order
	if err = cursor.All(ctx, &orders); err != nil {
		log.Printf("Error decoding all orders from cursor: %v", err)
		return nil, errors.New("failed to process all order data")
	}

	var orderResponses []OrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, *orderToResponse(&o))
	}

	return orderResponses, nil
}

// UpdateOrderStatus updates the status of an order.
func (s *service) UpdateOrderStatus(ctx context.Context, orderID string, req *UpdateOrderStatusRequest) (*OrderResponse, error) {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID format")
	}

	update := bson.M{
		"status":    req.Status,
		"updatedAt": time.Now(),
	}

	result := s.ordersCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		log.Printf("Error updating order status: %v", result.Err())
		return nil, errors.New("failed to update order status")
	}

	var updatedOrder Order
	if err := result.Decode(&updatedOrder); err != nil {
		log.Printf("Error decoding updated order: %v", err)
		return nil, errors.New("failed to decode updated order data")
	}

	return orderToResponse(&updatedOrder), nil
}
