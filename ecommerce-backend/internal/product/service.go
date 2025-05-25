// internal/product/service.go
package product

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // For find options like limit, skip, sort

	"github.com/u-s1ddhar7h/ecommerce-backend/internal/database" // Import your database package
)

// ProductService defines the interface for product operations.
type ProductService interface {
	CreateProduct(ctx context.Context, req *ProductCreateRequest) (*ProductResponse, error)
	GetProductByID(ctx context.Context, id string) (*ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]ProductResponse, error) // For now, no filters/pagination
	UpdateProduct(ctx context.Context, id string, req *ProductUpdateRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
	GetProductForOrder(ctx context.Context, id string) (*Product, error) // Internal use for order processing
}

// service implements ProductService.
type service struct {
	productsCollection *mongo.Collection
}

// NewProductService creates a new product service.
func NewProductService() ProductService {
	return &service{
		productsCollection: database.GetCollection("products"), // Get the 'products' collection
	}
}

// productToResponse converts a Product model to a ProductResponse.
func productToResponse(p *Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		SKU:         p.SKU,
		CategoryID:  p.CategoryID.Hex(),
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// CreateProduct handles the creation of a new product.
func (s *service) CreateProduct(ctx context.Context, req *ProductCreateRequest) (*ProductResponse, error) {
	// Convert CategoryID string to ObjectID
	categoryObjectID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	// Check if a product with the same SKU already exists
	var existingProduct Product
	err = s.productsCollection.FindOne(ctx, bson.M{"sku": req.SKU}).Decode(&existingProduct)
	if err == nil {
		return nil, errors.New("product with this SKU already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking for existing product SKU: %v", err)
		return nil, errors.New("database error during SKU check")
	}

	now := time.Now()
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		CategoryID:  categoryObjectID,
		Stock:       req.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result, err := s.productsCollection.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error inserting new product: %v", err)
		return nil, errors.New("failed to create product")
	}

	product.ID = result.InsertedID.(primitive.ObjectID)

	return productToResponse(product), nil
}

// GetProductByID retrieves a product by its ID.
func (s *service) GetProductByID(ctx context.Context, id string) (*ProductResponse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	var product Product
	err = s.productsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		log.Printf("Error finding product by ID: %v", err)
		return nil, errors.New("database error retrieving product")
	}

	return productToResponse(&product), nil
}

// GetAllProducts retrieves all products.
// In a real application, you'd add pagination and filtering here.
func (s *service) GetAllProducts(ctx context.Context) ([]ProductResponse, error) {
	cursor, err := s.productsCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})) // Sort by creation date descending
	if err != nil {
		log.Printf("Error finding all products: %v", err)
		return nil, errors.New("failed to retrieve products")
	}
	defer cursor.Close(ctx) // Ensure the cursor is closed

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("Error decoding products from cursor: %v", err)
		return nil, errors.New("failed to process product data")
	}

	var productResponses []ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *productToResponse(&p))
	}

	return productResponses, nil
}

// UpdateProduct updates an existing product.
func (s *service) UpdateProduct(ctx context.Context, id string, req *ProductUpdateRequest) (*ProductResponse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	update := bson.M{}
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Price != nil {
		update["price"] = *req.Price
	}
	if req.SKU != nil {
		update["sku"] = *req.SKU
	}
	if req.CategoryID != nil {
		categoryObjectID, err := primitive.ObjectIDFromHex(*req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category ID format for update")
		}
		update["categoryID"] = categoryObjectID
	}
	if req.Stock != nil {
		update["stock"] = *req.Stock
	}

	if len(update) == 0 {
		return nil, errors.New("no fields provided for update")
	}

	update["updatedAt"] = time.Now() // Update the timestamp on any change

	// Use $set to apply the updates
	result := s.productsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After), // Return the updated document
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		log.Printf("Error updating product: %v", result.Err())
		return nil, errors.New("failed to update product")
	}

	var updatedProduct Product
	if err := result.Decode(&updatedProduct); err != nil {
		log.Printf("Error decoding updated product: %v", err)
		return nil, errors.New("failed to decode updated product data")
	}

	return productToResponse(&updatedProduct), nil
}

// DeleteProduct deletes a product by its ID.
func (s *service) DeleteProduct(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}

	res, err := s.productsCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return errors.New("failed to delete product")
	}
	if res.DeletedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}

// GetProductForOrder is an internal helper to retrieve product details needed for order processing.
// It returns the full Product struct, not just the response version, as order logic needs stock.
func (s *service) GetProductForOrder(ctx context.Context, id string) (*Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	var product Product
	err = s.productsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		log.Printf("Error finding product for order: %v", err)
		return nil, errors.New("database error retrieving product for order")
	}
	return &product, nil
}
