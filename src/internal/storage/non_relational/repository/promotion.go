package nonrelational

import (
	"context"
	"errors"
	"fmt"
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
	discountsCollection := r.db.Collection("discounts")
	financingsCollection := r.db.Collection("financings")

	logger.Info("Finding promotions for store with CUIT %s between %v and %v in non-relational repository.", cuit, startDate, endDate)

	var promotionsDiscount []models.Discount
	var promotionsFinancing []models.Financing

	// Build the query
	filter := bson.M{
		"promotion_entity.cuit_store": cuit,
		"$and": []bson.M{
			{"promotion_entity.validity_start_date": bson.M{"$gte": startDate, "$lte": endDate}},
			{"promotion_entity.validity_end_date": bson.M{"$gte": startDate, "$lte": endDate}},
		},
	}

	// Query discounts
	cursor, err := discountsCollection.Find(context.TODO(), filter)
	if err != nil {
		logger.Error("Error finding DiscountEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}
	defer cursor.Close(context.TODO())

	var discounts []entities.DiscountEntityNonSQL
	if err := cursor.All(context.TODO(), &discounts); err != nil {
		return nil, nil, err
	}

	logger.Info("Found %d discounts for store with CUIT %s between %v and %v", len(discounts), cuit, startDate, endDate)

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

	logger.Info("Found %d financings for store with CUIT %s between %v and %v", len(financings), cuit, startDate, endDate)

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

	logger.Info("Finding most used promotion on Non-Relational database")

	// Collection references
	monthlyPurchasesCollection := r.db.Collection("purchase_monthly_payments")
	singlePurchasesCollection := r.db.Collection("purchase_single_payments")

	// Aggregation pipeline to count occurrences of payment_voucher in single payments
	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$purchase.payment_voucher"},
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

	// Aggregate for single purchases
	cursorSingle, err := singlePurchasesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		logger.Info("Error aggregating single purchases for most used promotion: %v", err)
		return nil, err
	}
	defer cursorSingle.Close(context.TODO())

	// Aggregate for monthly purchases
	cursorMonthly, err := monthlyPurchasesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		logger.Info("Error aggregating monthly purchases for most used promotion: %v", err)
		return nil, err
	}
	defer cursorMonthly.Close(context.TODO())

	// Decode results
	var singleResult, monthlyResult struct {
		PaymentVoucher    string `bson:"_id"`
		TotalRepeticiones int    `bson:"total_repeticiones"`
	}

	hasSingle := cursorSingle.Next(context.TODO())
	if hasSingle {
		if err := cursorSingle.Decode(&singleResult); err != nil {
			return nil, fmt.Errorf("error decoding single purchases: %w", err)
		}
	}

	hasMonthly := cursorMonthly.Next(context.TODO())
	if hasMonthly {
		if err := cursorMonthly.Decode(&monthlyResult); err != nil {
			return nil, fmt.Errorf("error decoding monthly purchases: %w", err)
		}
	}

	// Determine the most used promotion between single and monthly
	var mostUsedPromotion string
	if hasSingle && hasMonthly {
		if singleResult.TotalRepeticiones >= monthlyResult.TotalRepeticiones {
			logger.Info("Most used promotion found for single payments vs monthly payments: %s", singleResult.PaymentVoucher)
			mostUsedPromotion = singleResult.PaymentVoucher
		} else {
			logger.Info("Most used promotion found for monthly payments vs single payments: %s", monthlyResult.PaymentVoucher)
			mostUsedPromotion = monthlyResult.PaymentVoucher
		}
	} else if hasSingle {
		logger.Info("Most used promotion found for single payments: %s", singleResult.PaymentVoucher)
		mostUsedPromotion = singleResult.PaymentVoucher
	} else if hasMonthly {
		logger.Info("Most used promotion found for monthly payments: %s", monthlyResult.PaymentVoucher)
		mostUsedPromotion = monthlyResult.PaymentVoucher
	} else {
		return nil, errors.New("no promotions found")
	}

	// Fetch the promotion details using the most used promotion code
	return r.findPromotionByCode(mostUsedPromotion)
}

func (r *PromotionRepositoryMongo) findPromotionByCode(code string) (interface{}, error) {
	logger.Info("Finding promotion with code %s", code)

	discountsCollection := r.db.Collection("discounts")
	financingsCollection := r.db.Collection("financings")

	// Try to find in discount promotions
	var discount entities.DiscountEntityNonSQL
	if err := discountsCollection.FindOne(context.TODO(), bson.M{"promotion_entity.code": code}).Decode(&discount); err == nil {
		logger.Info("Found discount promotion with code %s", code)
		return discount, nil
	}

	// Try to find in financing promotions
	var financing entities.FinancingEntityNonSQL
	if err := financingsCollection.FindOne(context.TODO(), bson.M{"promotion_entity.code": code}).Decode(&financing); err == nil {
		logger.Info("Found financing promotion with code %s", code)
		return financing, nil
	}

	return nil, errors.New("No discount or financing promotion found with code " + code)
}
