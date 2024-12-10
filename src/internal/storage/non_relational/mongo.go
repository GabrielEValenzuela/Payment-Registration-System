package nonrelational

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewMongoDB establishes a new connection to the MongoDB database and initializes the collections.
func NewMongoDB(URI string, Database string, CleanDB bool) (*mongo.Database, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the MongoDB server to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Connected to MongoDB successfully.")

	// Get the specified database
	db := client.Database(Database)

	// Initialize the database schema
	if err := initMongoDB(ctx, db, CleanDB); err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB database: %w", err)
	}

	return db, nil
}

// CloseMongoDB gracefully closes the MongoDB connection.
func CloseMongoDB(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to close MongoDB connection: %w", err)
	}

	logger.Info("MongoDB connection closed successfully.")
	return nil
}

// initMongoDB initializes the MongoDB database schema.
// If `cleanDB` is true, it will drop existing collections.
func initMongoDB(ctx context.Context, db *mongo.Database, cleanDB bool) error {
	if cleanDB {
		logger.Info("Cleaning the database: dropping existing collections...")
		collections, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("failed to list collections: %w", err)
		}

		for _, collection := range collections {
			if err := db.Collection(collection).Drop(ctx); err != nil {
				return fmt.Errorf("failed to drop collection %s: %w", collection, err)
			}
			logger.Info("Dropped collection: %s", collection)
		}
	}

	// Create indexes or seed data here if needed
	logger.Info("MongoDB schema initialized successfully.")
	return nil
}
