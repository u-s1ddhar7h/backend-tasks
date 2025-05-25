// internal/database/mongodb.go
package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient is the global MongoDB client instance.
var MongoClient *mongo.Client

// ConnectDB establishes a connection to MongoDB.
func ConnectDB(mongoURI string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Create a context with a timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is cancelled to release resources

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	MongoClient = client // Store the connected client globally
}

// GetCollection returns a specific MongoDB collection.
// You might want to pass the database name dynamically or define it in config.
func GetCollection(collectionName string) *mongo.Collection {
	// "ecommerce" is the name of your database. Make sure it matches your MONGO_URI if different.
	return MongoClient.Database("ecommerce").Collection(collectionName)
}
