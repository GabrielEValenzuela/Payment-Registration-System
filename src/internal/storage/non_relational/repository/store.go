package nonrelational

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type StoreRepositoryMongo struct {
	db *mongo.Database
}

// NewStoreNonRelationalRepository creates a new instance of StoreRepositoryMongo
func NewStoreNonRelationalRepository(db *mongo.Database) *StoreRepositoryMongo {
	return &StoreRepositoryMongo{db: db}
}

// GetStoreWithHighestRevenueByMonth retrieves the store with the highest revenue for a specific month and year.
func (r *StoreRepositoryMongo) GetStoreWithHighestRevenueByMonth(month int, year int) (models.StoreDTO, error) {
	purchasesCollection := r.db.Collection("purchases")

	// Define the start and end dates for the given month and year
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// Aggregation pipeline
	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			}},
		},
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "store", Value: "$store"},
					{Key: "cuit_store", Value: "$cuit_store"},
				}},
				{Key: "total_amount", Value: bson.M{"$sum": "$final_amount"}},
			}},
		},
		{
			{Key: "$sort", Value: bson.M{
				"total_amount": -1,
			}},
		},
		{
			{Key: "$limit", Value: 1},
		},
	}

	cursor, err := purchasesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return models.StoreDTO{}, fmt.Errorf("error retrieving store with highest revenue: %v", err)
	}
	defer cursor.Close(context.TODO())

	var result struct {
		ID struct {
			Name string `bson:"store"`
			Cuit string `bson:"cuit_store"`
		} `bson:"_id"`
		TotalAmount float64 `bson:"total_amount"`
	}

	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return models.StoreDTO{}, fmt.Errorf("error decoding aggregation result: %v", err)
		}
	} else {
		return models.StoreDTO{}, fmt.Errorf("no results found for month %d, year %d", month, year)
	}

	// Map the result to the StoreDTO model
	storeDTO := models.StoreDTO{
		Name: result.ID.Name,
		Cuit: result.ID.Cuit,
	}

	return storeDTO, nil
}
