package nonrelational

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CardRepositoryMongo struct {
	db *mongo.Database
}

// NewCardRepository creates a new instance of CardRepositoryMongo
func NewCardRepositoryNonRelational(db *mongo.Database) *CardRepositoryMongo {
	return &CardRepositoryMongo{db: db}
}

func (r *CardRepositoryMongo) GetPaymentSummary(cardNumber string, month int, year int) (*models.PaymentSummary, error) {
	collection := r.db.Collection("cards")

	// Define the date range for the given month and year
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // One month later

	// Find the card by its number
	var card entities.CardEntityNonSQL
	if err := collection.FindOne(context.TODO(), bson.M{"number": cardNumber}).Decode(&card); err != nil {
		return nil, fmt.Errorf("could not find card with number %s: %v", cardNumber, err)
	}

	// Filter purchases within the date range
	var totalPrice float64
	for _, purchase := range card.PurchaseSinglePayments {
		if purchase.PurchaseEntity.CreatedAt.After(startDate) && purchase.PurchaseEntity.CreatedAt.Before(endDate) {
			totalPrice += purchase.PurchaseEntity.Amount
		}
	}
	for _, purchase := range card.PurchaseMonthlyPayments {
		if purchase.PurchaseEntity.CreatedAt.After(startDate) && purchase.PurchaseEntity.CreatedAt.Before(endDate) {
			totalPrice += purchase.PurchaseEntity.Amount
		}
	}

	// Define expiration dates
	firstExpiration := time.Now().AddDate(0, 0, 15)       // 15 days from today
	secondExpiration := firstExpiration.AddDate(0, 0, 10) // 10 days from today

	// Generate a unique code for the Payment Summary
	code := fmt.Sprintf("SUMMARY-%d-%d", year, month)

	// Create the PaymentSummary object
	paymentSummary := models.PaymentSummary{
		Code:                code,
		Month:               month,
		Year:                year,
		FirstExpiration:     firstExpiration,
		SecondExpiration:    secondExpiration,
		SurchargePercentage: 5.0,        // Example: a 5% surcharge
		TotalPrice:          totalPrice, // Total of all purchases
		Card:                *entities.ToCard(&card),
	}

	return &paymentSummary, nil
}

func (r *CardRepositoryMongo) GetCardsExpiringInNext30Days(day int, month int, year int) (*[]models.Card, error) {
	collection := r.db.Collection("cards")

	// Define the date range
	startDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 30)

	// Query for cards with expiration dates in the next 30 days
	cursor, err := collection.Find(context.TODO(), bson.M{
		"expiration_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error querying cards expiring in the next 30 days: %v", err)
	}
	defer cursor.Close(context.TODO())

	var cards []models.Card
	for cursor.Next(context.TODO()) {
		var card entities.CardEntityNonSQL
		if err := cursor.Decode(&card); err != nil {
			return nil, err
		}
		cards = append(cards, *entities.ToCard(&card))
	}

	return &cards, nil
}

func (r *CardRepositoryMongo) GetPurchaseSingle(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseSinglePayment, error) {
	collection := r.db.Collection("purchases")

	// Query for a single-payment purchase by CUIT, final amount, and payment voucher
	var purchase entities.PurchaseSinglePaymentEntityNonSQL
	if err := collection.FindOne(context.TODO(), bson.M{
		"cuit_store":      cuit,
		"final_amount":    finalAmount,
		"payment_voucher": paymentVoucher,
	}).Decode(&purchase); err != nil {
		return nil, fmt.Errorf("error finding single-payment purchase: %v", err)
	}

	return entities.ToPurchaseSinglePaymentNonSQL(&purchase), nil
}

func (r *CardRepositoryMongo) GetPurchaseMonthly(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseMonthlyPayment, error) {
	collection := r.db.Collection("purchases")

	// Query for a monthly-payment purchase by CUIT, final amount, and payment voucher
	var purchase entities.PurchaseMonthlyPaymentsEntityNonSQL
	if err := collection.FindOne(context.TODO(), bson.M{
		"cuit_store":      cuit,
		"final_amount":    finalAmount,
		"payment_voucher": paymentVoucher,
	}).Decode(&purchase); err != nil {
		return nil, fmt.Errorf("error finding monthly-payment purchase: %v", err)
	}

	return entities.ToPurchaseMonthlyPaymentsNonSQL(&purchase), nil
}

func (r *CardRepositoryMongo) GetTop10CardsByPurchases() (*[]models.Card, error) {
	collection := r.db.Collection("cards")

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "number", Value: 1},
				{Key: "purchase_count", Value: bson.D{
					{Key: "$add", Value: bson.A{
						bson.D{{Key: "$size", Value: "$purchase_single_payments"}},
						bson.D{{Key: "$size", Value: "$purchase_monthly_payments"}},
					}},
				}},
			}},
		},
		{{
			Key: "$sort", Value: bson.D{
				{Key: "purchase_count", Value: -1},
			},
		}},
		{{
			Key: "$limit", Value: 10,
		}},
	}

	// Aggregate the top 10 cards by total number of purchases
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating top 10 cards: %w", err)
	}
	defer cursor.Close(context.TODO())

	// Decode results into the models.Card slice
	var cards []models.Card
	for cursor.Next(context.TODO()) {
		var card entities.CardEntityNonSQL
		if err := cursor.Decode(&card); err != nil {
			return nil, fmt.Errorf("error decoding card: %w", err)
		}
		cards = append(cards, *entities.ToCard(&card))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return &cards, nil
}
