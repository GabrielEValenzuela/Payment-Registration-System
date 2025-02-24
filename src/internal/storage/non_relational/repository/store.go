package nonrelational

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
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

	logger.Info("Retrieving store with highest revenue for month %d, year %d", month, year)

	singlePurchasesCollection := r.db.Collection("purchase_single_payments")
	monthlyPurchasesCollection := r.db.Collection("purchase_monthly_payments")

	// Define the start and end dates for the given month and year
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{
				"purchase.created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			}},
		},
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "store", Value: "$purchase.store"},
					{Key: "cuit_store", Value: "$purchase.cuit_store"},
				}},
				{Key: "total_amount", Value: bson.M{"$sum": "$purchase.final_amount"}},
			}},
		},
		{
			{Key: "$sort", Value: bson.M{"total_amount": -1}},
		},
		{
			{Key: "$limit", Value: 1},
		},
	}

	// Process single payments
	cursorSingle, err := singlePurchasesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return models.StoreDTO{}, fmt.Errorf("error retrieving store with highest revenue from single purchases: %w", err)
	}
	defer cursorSingle.Close(context.TODO())

	// Process monthly payments
	cursorMonthly, err := monthlyPurchasesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return models.StoreDTO{}, fmt.Errorf("error retrieving store with highest revenue from monthly purchases: %w", err)
	}
	defer cursorMonthly.Close(context.TODO())

	// Decode results
	var singleResult, monthlyResult struct {
		ID struct {
			Name string `bson:"store"`
			Cuit string `bson:"cuit_store"`
		} `bson:"_id"`
		TotalAmount float64 `bson:"total_amount"`
	}

	hasSingle := cursorSingle.Next(context.TODO())
	if hasSingle {
		if err := cursorSingle.Decode(&singleResult); err != nil {
			return models.StoreDTO{}, fmt.Errorf("error decoding single purchases: %w", err)
		}
		logger.Info("Single payment result: %v", singleResult)
	}

	hasMonthly := cursorMonthly.Next(context.TODO())
	if hasMonthly {
		if err := cursorMonthly.Decode(&monthlyResult); err != nil {
			return models.StoreDTO{}, fmt.Errorf("error decoding monthly purchases: %w", err)
		}
		logger.Info("Monthly payment result: %v", monthlyResult)
	}

	// Determine which store has the highest revenue
	var bestStore struct {
		Name        string
		Cuit        string
		TotalAmount float64
	}

	if hasSingle && hasMonthly {
		if singleResult.TotalAmount >= monthlyResult.TotalAmount {
			bestStore = struct {
				Name        string
				Cuit        string
				TotalAmount float64
			}{singleResult.ID.Name, singleResult.ID.Cuit, singleResult.TotalAmount}
			logger.Info("Single payment result is higher: %v", bestStore)
		} else {
			bestStore = struct {
				Name        string
				Cuit        string
				TotalAmount float64
			}{monthlyResult.ID.Name, monthlyResult.ID.Cuit, monthlyResult.TotalAmount}
			logger.Info("Monthly payment result is higher: %v", bestStore)
		}
	} else if hasSingle {
		bestStore = struct {
			Name        string
			Cuit        string
			TotalAmount float64
		}{singleResult.ID.Name, singleResult.ID.Cuit, singleResult.TotalAmount}
	} else if hasMonthly {
		bestStore = struct {
			Name        string
			Cuit        string
			TotalAmount float64
		}{monthlyResult.ID.Name, monthlyResult.ID.Cuit, monthlyResult.TotalAmount}
	} else {
		return models.StoreDTO{}, fmt.Errorf("no results found for month %d, year %d", month, year)
	}

	// Map the result to the StoreDTO model
	storeDTO := models.StoreDTO{
		Name: bestStore.Name,
		Cuit: bestStore.Cuit,
	}

	return storeDTO, nil
}
