/*
 * Payment Registration System - Non-Relational Repository
 * --------------------------------------------------------
 * This file contains the implementation of the BankRepositoryMongo struct.
 * The struct implements the IBankStorage interface and provides methods to interact with the MongoDB database.
 *
 * Created: Dec. 10, 2024
 * License: GNU General Public License v3.0
 */

package nonrelational

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type BankRepositoryMongo struct {
	db *mongo.Database
}

func NewBankNonRelationalRepository(db *mongo.Database) storage.IBankStorage {
	return &BankRepositoryMongo{db: db}
}

// AddFinancingPromotionToBank adds a financing promotion to a bank.
func (r *BankRepositoryMongo) AddFinancingPromotionToBank(promotionFinancing models.Financing) error {
	ctx := context.Background()

	bankCollection := r.db.Collection("banks")
	var bank entities.BankEntityNonSQL
	if err := bankCollection.FindOne(ctx, bson.M{"cuit": promotionFinancing.Promotion.Bank.Cuit}).Decode(&bank); err != nil {
		return fmt.Errorf("could not find bank with 'cuit' %s: %w", promotionFinancing.Promotion.Bank.Cuit, err)
	}

	// Fetched bank
	logger.Info("Information for bank 'cuit' %s: %v", promotionFinancing.Promotion.Bank.Cuit, bank)

	financingEntity := entities.FinancingEntityNonSQL{
		PromotionEntity: entities.PromotionEntityNonSQL{
			Code:              promotionFinancing.Promotion.Code,
			PromotionTitle:    promotionFinancing.Promotion.PromotionTitle,
			NameStore:         promotionFinancing.Promotion.NameStore,
			CuitStore:         promotionFinancing.Promotion.CuitStore,
			ValidityStartDate: promotionFinancing.Promotion.ValidityStartDate.Time,
			ValidityEndDate:   promotionFinancing.Promotion.ValidityEndDate.Time,
			Comments:          promotionFinancing.Promotion.Comments,
		},
		NumberOfQuotas: promotionFinancing.NumberOfQuotas,
		Interest:       promotionFinancing.Interest,
		IsDeleted:      false,
		BankID:         bank.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.db.Collection("financings").InsertOne(ctx, financingEntity)
	if err != nil {
		logger.Error("Failed to add financing promotion %v", err)
		return fmt.Errorf("could not add financing promotion: %w", err)
	}

	logger.Info("Successfully added financing promotion %s to bank %s in NoSQL storage",
		promotionFinancing.Promotion.Code, promotionFinancing.Promotion.Bank.Cuit)

	return nil

}

// ExtendPromotionValidity extends the validity of a promotion.
func (r *BankRepositoryMongo) ExtendPromotionValidity(code string, newDate time.Time) error {

	ctx := context.Background()

	// Update the promotion
	filter := bson.M{"code": code}
	update := bson.M{"$set": bson.M{"validity_end_date": newDate}}
	result, err := r.db.Collection("promotions").UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update promotion %s: %v", code, err)
		return err
	}

	// If no document was modified, return an explicit error
	if result.ModifiedCount == 0 {
		logger.Error("No promotion found with code %s or validity date was already up to date", code)
		return err
	}

	logger.Info("Successfully extended promotion %s validity to %s", code, newDate)
	return nil
}

// DeletePromotion marks a promotion as deleted.
func (r *BankRepositoryMongo) DeletePromotion(code string) error {

	ctx := context.Background()

	// Mark the promotion as deleted
	filter := bson.M{"promotion_entity.code": code}
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
	promotionCollection := r.db.Collection("financings")
	var promotion bson.M
	if err := promotionCollection.FindOne(ctx, bson.M{"promotion_entity.code": code}).Decode(&promotion); err != nil {
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
	promotionCollection := r.db.Collection("discounts")
	var promotion bson.M
	if err := promotionCollection.FindOne(ctx, bson.M{"promotion_entity.code": code}).Decode(&promotion); err != nil {
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
	result, err := r.db.Collection("discounts").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update promotion %s: %w", code, err)
	}

	// If no document was modified, return an explicit error
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no promotion found with code %s or validity date was already up to date", code)
	}

	logger.Info("Successfully extended promotion %s validity to %s", code, newDate)
	return nil
}

// ExtendFinancingPromotionValidity extends the validity of a promotion.
func (r *BankRepositoryMongo) ExtendFinancingPromotionValidity(code string, newDate time.Time) error {

	ctx := context.Background()

	// Update the promotion
	filter := bson.M{"promotion_entity.code": code}
	update := bson.M{"$set": bson.M{"validity_end_date": newDate}}
	result, err := r.db.Collection("financings").UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to update promotion %s: %v", code, err)
		return err
	}

	// If no document was modified, return an explicit error
	if result.ModifiedCount == 0 {
		logger.Error("No promotion found with code %s or validity date was already up to date", code)
		return err
	}

	logger.Info("Successfully extended financing promotion %s validity. New end date: '%s'", code, newDate)
	return nil
}

// GetBankCustomerCounts retrieves the number of customers for each bank.
func (r *BankRepositoryMongo) GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error) {
	ctx := context.Background()

	// Reference the banks collection
	bankCollection := r.db.Collection("banks")
	pipeline := []bson.M{
		// Lookup to join BANKS with CUSTOMERS_BANKS
		{
			"$lookup": bson.M{
				"from":         "customers_banks",
				"localField":   "_id",     // BANKS _id field
				"foreignField": "bank_id", // CUSTOMERS_BANKS bank_id field
				"as":           "customer_relations",
			},
		},
		// Add a field to count the number of customers
		{
			"$addFields": bson.M{
				"customer_count": bson.M{"$size": "$customer_relations"},
			},
		},
		// Project the desired fields
		{
			"$project": bson.M{
				"_id":            0,
				"bank_cuit":      "$cuit",
				"bank_name":      "$name",
				"customer_count": 1,
			},
		},
	}

	// Execute the aggregation pipeline
	cursor, err := bankCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve customer counts: %w", err)
	}

	// Decode the results into a map first to normalize the data
	var rawResults []bson.M
	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, fmt.Errorf("could not decode results: %w", err)
	}

	// Convert raw results into DTOs
	var results []models.BankCustomerCountDTO
	for _, raw := range rawResults {
		customerCount, ok := raw["customer_count"].(int32) // MongoDB often stores counts as int32
		if !ok {
			return nil, fmt.Errorf("invalid customer_count format")
		}

		results = append(results, models.BankCustomerCountDTO{
			BankCuit:      raw["bank_cuit"].(string),
			BankName:      raw["bank_name"].(string),
			CustomerCount: int(customerCount),
		})
	}

	return results, nil
}
