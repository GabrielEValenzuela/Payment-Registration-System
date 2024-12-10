package nonrelational

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type BankRepositoryMongo struct {
	db *mongo.Database
}

func NewBankNonRelationalRepository(db *mongo.Database) storage.IStorage {
	return &BankRepositoryMongo{db: db}
}

// AddFinancingPromotionToBank adds a financing promotion to a bank.
func (r *BankRepositoryMongo) AddFinancingPromotionToBank(promotionFinancing models.Financing) error {

	ctx := context.Background()

	// Find the bank by CUIT
	bankCollection := r.db.Collection("banks")
	var bank bson.M
	if err := bankCollection.FindOne(ctx, bson.M{"cuit": promotionFinancing.Bank.Cuit}).Decode(&bank); err != nil {
		return fmt.Errorf("could not find bank with CUIT %s: %w", promotionFinancing.Bank.Cuit, err)
	}

	// Add the promotion to the bank
	promotion := bson.M{
		"code":                promotionFinancing.Code,
		"promotion_title":     promotionFinancing.PromotionTitle,
		"validity_start_date": promotionFinancing.ValidityStartDate,
		"validity_end_date":   promotionFinancing.ValidityEndDate,
		"is_deleted":          false,
		"bank_id":             bank["_id"],
	}
	_, err := r.db.Collection("financing_promotions").InsertOne(ctx, promotion)
	if err != nil {
		logger.Error("Failed to add financing promotion to bank: %v", zap.Error(err))
		return fmt.Errorf("could not add financing promotion: %w", err)
	}

	logger.Info("Successfully added financing promotion %s to bank %s", promotionFinancing.Code, promotionFinancing.Bank.Cuit)
	return nil
}

// ExtendPromotionValidity extends the validity of a promotion.
func (r *BankRepositoryMongo) ExtendPromotionValidity(code string, newDate time.Time) error {

	ctx := context.Background()

	// Update the promotion
	filter := bson.M{"code": code}
	update := bson.M{"$set": bson.M{"validity_end_date": newDate}}
	result, err := r.db.Collection("promotions").UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("could not extend promotion validity: %w", err)
	}

	logger.Info("Successfully extended promotion %s validity to %s", code, newDate)
	return nil
}

// DeletePromotion marks a promotion as deleted.
func (r *BankRepositoryMongo) DeletePromotion(code string) error {

	ctx := context.Background()

	// Mark the promotion as deleted
	filter := bson.M{"code": code}
	update := bson.M{"$set": bson.M{"is_deleted": true}}
	result, err := r.db.Collection("promotions").UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("could not delete promotion: %w", err)
	}

	logger.Info("Successfully marked promotion %s as deleted", code)
	return nil
}

// DeleteFinancingPromotion deletes a financing promotion from the bank.
func (r *BankRepositoryMongo) DeleteFinancingPromotion(code string) error {

	ctx := context.Background()

	// Find the promotion by code
	promotionCollection := r.db.Collection("financing_promotions")
	var promotion bson.M
	if err := promotionCollection.FindOne(ctx, bson.M{"code": code}).Decode(&promotion); err != nil {
		return fmt.Errorf("could not find promotion with code %s: %w", code, err)
	}

	// Delete the promotion
	_, err := promotionCollection.DeleteOne(ctx, bson.M{"code": code})
	if err != nil {
		logger.Error("Failed to delete financing promotion: %v", zap.Error(err))
		return fmt.Errorf("could not delete financing promotion: %w", err)
	}

	logger.Info("Successfully deleted financing promotion %s", code)
	return nil
}

// DeleteDiscountPromotion deletes a discount promotion from the bank.
func (r *BankRepositoryMongo) DeleteDiscountPromotion(code string) error {

	ctx := context.Background()

	// Find the promotion by code
	promotionCollection := r.db.Collection("discount_promotions")
	var promotion bson.M
	if err := promotionCollection.FindOne(ctx, bson.M{"code": code}).Decode(&promotion); err != nil {
		return fmt.Errorf("could not find promotion with code %s: %w", code, err)
	}

	// Delete the promotion
	_, err := promotionCollection.DeleteOne(ctx, bson.M{"code": code})
	if err != nil {
		logger.Error("Failed to delete discount promotion: %v", zap.Error(err))
		return fmt.Errorf("could not delete discount promotion: %w", err)
	}

	logger.Info("Successfully deleted discount promotion %s", code)
	return nil
}

// ExtendDiscountPromotionValidity extends the validity of a promotion.
func (r *BankRepositoryMongo) ExtendDiscountPromotionValidity(code string, newDate time.Time) error {

	ctx := context.Background()

	// Update the promotion
	filter := bson.M{"code": code}
	update := bson.M{"$set": bson.M{"validity_end_date": newDate}}
	result, err := r.db.Collection("discount_promotions").UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("could not extend promotion validity: %w", err)
	}

	logger.Info("Successfully extended promotion %s validity to %s", code, newDate)
	return nil
}

// ExtendFinancingPromotionValidity extends the validity of a promotion.
func (r *BankRepositoryMongo) ExtendFinancingPromotionValidity(code string, newDate time.Time) error {

	ctx := context.Background()

	// Update the promotion
	filter := bson.M{"code": code}
	update := bson.M{"$set": bson.M{"validity_end_date": newDate}}
	result, err := r.db.Collection("financing_promotions").UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("could not extend promotion validity: %w", err)
	}

	logger.Info("Successfully extended promotion %s validity to %s", code, newDate)
	return nil
}

// GetBankCustomerCounts retrieves the number of customers for each bank.
func (r *BankRepositoryMongo) GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error) {

	ctx := context.Background()

	customerCollection := r.db.Collection("customers")
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "banks",
				"localField":   "bank_id",
				"foreignField": "_id",
				"as":           "bank_info",
			},
		},
		{
			"$group": bson.M{
				"_id":            "$bank_info",
				"customer_count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := customerCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve customer counts: %w", err)
	}

	var results []models.BankCustomerCountDTO
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("could not decode results: %w", err)
	}

	return results, nil
}
