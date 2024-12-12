package nonrelational

import (
	"context"
	"errors"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PromotionRepositoryMongo struct {
	db *mongo.Database
}

// NewPromotionRepository creates a new instance of PromotionRepositoryMongo
func NewPromotionNonRelationalRepository(db *mongo.Database) storage.IPromotionStorage {
	return &PromotionRepositoryMongo{db: db}
}

// GetAvailablePromotionsByStoreAndDateRange retrieves available promotions for a store within a date range.
func (r *PromotionRepositoryMongo) GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate, endDate time.Time) (*[]models.Financing, *[]models.Discount, error) {
	discountsCollection := r.db.Collection("discount_promotions")
	financingsCollection := r.db.Collection("financing_promotions")

	var promotionsDiscount []models.Discount
	var promotionsFinancing []models.Financing

	// Build the query
	filter := bson.M{
		"cuit_store": cuit,
		"is_deleted": false,
		"$or": []bson.M{
			{"validity_start_date": bson.M{"$gte": startDate, "$lte": endDate}},
			{"validity_end_date": bson.M{"$gte": startDate, "$lte": endDate}},
		},
	}

	// Query discounts
	cursor, err := discountsCollection.Find(context.TODO(), filter)
	if err != nil {
		logger.Info("Error finding DiscountEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}
	defer cursor.Close(context.TODO())

	var discounts []entities.DiscountEntityNonSQL
	if err := cursor.All(context.TODO(), &discounts); err != nil {
		return nil, nil, err
	}

	// Query financings
	cursor, err = financingsCollection.Find(context.TODO(), filter)
	if err != nil {
		logger.Info("Error finding FinancingEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}
	defer cursor.Close(context.TODO())

	var financings []entities.FinancingEntityNonSQL
	if err := cursor.All(context.TODO(), &financings); err != nil {
		return nil, nil, err
	}

	// Convert results to models
	for _, discount := range discounts {
		promotionsDiscount = append(promotionsDiscount, *entities.ToDiscountNonSQL(&discount))
	}
	for _, financing := range financings {
		promotionsFinancing = append(promotionsFinancing, *entities.ToFinancingNonSQL(&financing))
	}

	return &promotionsFinancing, &promotionsDiscount, nil
}

// GetMostUsedPromotion retrieves the most used promotion based on its usage.
func (r *PromotionRepositoryMongo) GetMostUsedPromotion() (interface{}, error) {
	purchasesCollection := r.db.Collection("purchases")

	// Aggregation pipeline to count occurrences of payment_voucher
	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$payment_voucher"},
				{Key: "total_repeticiones", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		},
		{
			{Key: "$sort", Value: bson.D{
				{Key: "total_repeticiones", Value: -1},
			}},
		},
		{
			{Key: "$limit", Value: 1},
		},
	}

	cursor, err := purchasesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		logger.Info("Error aggregating purchases for most used promotion: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		PaymentVoucher    string `bson:"_id"`
		TotalRepeticiones int    `bson:"total_repeticiones"`
	}
	if !cursor.Next(context.TODO()) {
		return nil, errors.New("no promotions found")
	}
	if err := cursor.Decode(&result); err != nil {
		return nil, err
	}

	// Find the promotion by payment voucher
	return r.findPromotionByCode(result.PaymentVoucher)
}

func (r *PromotionRepositoryMongo) findPromotionByCode(code string) (interface{}, error) {
	discountsCollection := r.db.Collection("discount_promotions")
	financingsCollection := r.db.Collection("financing_promotions")

	// Try to find in discount promotions
	var discount entities.DiscountEntityNonSQL
	if err := discountsCollection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&discount); err == nil {
		return entities.ToDiscountNonSQL(&discount), nil
	}

	// Try to find in financing promotions
	var financing entities.FinancingEntityNonSQL
	if err := financingsCollection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&financing); err == nil {
		return entities.ToFinancingNonSQL(&financing), nil
	}

	return nil, errors.New("promotion not found with the provided code")
}
