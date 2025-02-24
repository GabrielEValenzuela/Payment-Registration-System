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
)

type CardRepositoryMongo struct {
	db *mongo.Database
}

// NewCardRepository creates a new instance of CardRepositoryMongo
func NewCardNonRelationalRepository(db *mongo.Database) storage.ICardStorage {
	return &CardRepositoryMongo{db: db}
}

func (r *CardRepositoryMongo) GetPaymentSummary(cardNumber string, month int, year int) (*models.PaymentSummary, error) {
	collection := r.db.Collection("cards")

	// Define date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // Next month

	logger.Info("Fetching payment summary for card '%s' from %s to %s", cardNumber, startDate, endDate)

	// MongoDB Aggregation Pipeline
	pipeline := []bson.M{
		// 1) Match the specific card
		{"$match": bson.M{"number": cardNumber}},

		// 2) Lookup Bank details
		{"$lookup": bson.M{
			"from":         "banks",
			"localField":   "bank_cuit",
			"foreignField": "cuit",
			"as":           "bank",
		}},
		{"$unwind": bson.M{"path": "$bank", "preserveNullAndEmptyArrays": true}},

		// 3) Lookup Customer details
		{"$lookup": bson.M{
			"from":         "customers",
			"localField":   "customer_cuit",
			"foreignField": "cuit",
			"as":           "customer",
		}},
		{"$unwind": bson.M{"path": "$customer", "preserveNullAndEmptyArrays": true}},

		// 4) Lookup Single Payments *with* date filtering in the sub-pipeline
		{
			"$lookup": bson.M{
				"from": "purchase_single_payments",
				"let":  bson.M{"cNumber": "$number"},
				"pipeline": []bson.M{
					// Match same card_number
					{
						"$match": bson.M{
							"$expr": bson.M{"$eq": []interface{}{"$purchase.card_number", "$$cNumber"}},
						},
					},
					// Match created_at in range
					{
						"$match": bson.M{
							"purchase.created_at": bson.M{"$gte": startDate, "$lt": endDate},
						},
					},
				},
				"as": "single_payments",
			},
		},

		// 5) Lookup Monthly Payments *with* date filtering in the sub-pipeline
		{
			"$lookup": bson.M{
				"from": "purchase_monthly_payments",
				"let":  bson.M{"cNumber": "$number"},
				"pipeline": []bson.M{
					// Match same card_number
					{
						"$match": bson.M{
							"$expr": bson.M{"$eq": []interface{}{"$purchase.card_number", "$$cNumber"}},
						},
					},
					// Match created_at in range
					{
						"$match": bson.M{
							"purchase.created_at": bson.M{"$gte": startDate, "$lt": endDate},
						},
					},
				},
				"as": "monthly_payments",
			},
		},

		{
			"$project": bson.M{
				// Keep the root fields you care about
				"number":                  1,
				"ccv":                     1,
				"cardholder_name_in_card": 1,
				"bank_cuit":               1,
				"customer_cuit":           1,
				"created_at":              1,
				"updated_at":              1,

				// Bring along the arrays
				"single_payments":  1,
				"monthly_payments": 1,

				// Example total
				"total_price": bson.M{
					"$sum": []interface{}{
						// Sum of all single_payments.purchase.final_amount
						bson.M{"$sum": "$single_payments.purchase.final_amount"},
						// Sum of all monthly_payments.purchase.final_amount
						bson.M{"$sum": "$monthly_payments.purchase.final_amount"},
					},
				},
			},
		},
	}

	// Execute Aggregation Query
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching payment summary from MongoDB: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Decode Result
	if !cursor.Next(context.TODO()) {
		return nil, fmt.Errorf("no payment summary found for card %s", cardNumber)
	}

	var result entities.PaymentSummaryNoSQL
	if err := cursor.Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding payment summary: %v", err)
	}

	logger.Info("Payment generated: %+v", result)

	firstExpiration := time.Now().AddDate(0, 0, 15)       // 15 days from today
	secondExpiration := firstExpiration.AddDate(0, 0, 10) // 10 days later

	code := fmt.Sprintf("SUMMARY-%d-%d", year, month)

	paymentSummary := &models.PaymentSummary{
		Code:                code,
		Month:               month,
		Year:                year,
		FirstExpiration:     firstExpiration,
		SecondExpiration:    secondExpiration,
		SurchargePercentage: 0, // Assuming a default value, update as needed
		TotalPrice:          result.TotalPrice,
		SinglePayments:      *entities.ConvertPurchaseSinglePaymentListMongo(&result.SinglePayments),
		MonthlyPayments:     *entities.ConvertPurchaseMonthlyPaymentListMongo(&result.MonthlyPayments),
	}

	return paymentSummary, nil
}

func (r *CardRepositoryMongo) GetCardsExpiringInNext30Days(day int, month int, year int) (*[]models.Card, error) {
	collection := r.db.Collection("cards")

	// Define the date range
	startDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 30)

	logger.Info("Fetching cards expiring between %s and %s", startDate, endDate)

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

	logger.Info("Found %d cards expiring in the next 30 days", len(cards))

	return &cards, nil
}

func (r *CardRepositoryMongo) GetPurchaseSingle(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseSinglePayment, error) {
	collection := r.db.Collection("purchase_single_payments")

	logger.Info("Fetching single-payment purchase with CUIT: %s, final amount: %f, and payment voucher: %s", cuit, finalAmount, paymentVoucher)

	// Query for a single-payment purchase by nested fields
	var purchase entities.PurchaseSinglePaymentEntityNonSQL
	if err := collection.FindOne(
		context.TODO(),
		bson.M{
			"purchase.cuit_store":      cuit,
			"purchase.final_amount":    finalAmount,
			"purchase.payment_voucher": paymentVoucher,
		},
	).Decode(&purchase); err != nil {
		return nil, fmt.Errorf("error finding single-payment purchase: %v", err)
	}

	return entities.ToPurchaseSinglePaymentNonSQL(&purchase), nil
}

func (r *CardRepositoryMongo) GetPurchaseMonthly(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseMonthlyPayment, error) {
	collection := r.db.Collection("purchase_monthly_payments")

	logger.Info("Fetching monthly-payment purchase with CUIT: %s, final amount: %f, and payment voucher: %s",
		cuit, finalAmount, paymentVoucher)

	// Query for a single-payment purchase by nested fields
	var purchase entities.PurchaseMonthlyPaymentsEntityNonSQL
	if err := collection.FindOne(
		context.TODO(),
		bson.M{
			"purchase.cuit_store":      cuit,
			"purchase.final_amount":    finalAmount,
			"purchase.payment_voucher": paymentVoucher,
		},
	).Decode(&purchase); err != nil {
		return nil, fmt.Errorf("error finding single-payment purchase: %v", err)
	}

	return entities.ToPurchaseMonthlyPaymentsNonSQL(&purchase), nil
}

func (r *CardRepositoryMongo) GetTop10CardsByPurchases() (*[]models.Card, error) {
	collection := r.db.Collection("cards")

	logger.Info("Fetching top 10 cards by total number of purchases")

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// 1. Lookup single-payments by matching card "number"
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "purchase_single_payments"},
				{Key: "localField", Value: "number"},
				{Key: "foreignField", Value: "purchase.card_number"},
				{Key: "as", Value: "purchase_single_payments"},
			},
		}},
		// 2. Lookup monthly-payments by matching card "number"
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "purchase_monthly_payments"},
				{Key: "localField", Value: "number"},
				{Key: "foreignField", Value: "purchase.card_number"},
				{Key: "as", Value: "purchase_monthly_payments"},
			},
		}},
		// 3. Project fields + compute the total purchase_count
		bson.D{{
			Key: "$project",
			Value: bson.D{
				{Key: "id", Value: 1},
				{Key: "number", Value: 1},
				{Key: "ccv", Value: 1},
				{Key: "cardholder_name_in_card", Value: 1},
				{Key: "since", Value: 1},
				{Key: "expiration_date", Value: 1},
				{Key: "bank_cuit", Value: 1},
				{Key: "customer_cuit", Value: 1},
				{Key: "purchase_single_payments", Value: 1},
				{Key: "purchase_monthly_payments", Value: 1},
				{
					Key: "purchase_count",
					Value: bson.D{
						{
							Key: "$add",
							Value: bson.A{
								bson.D{{Key: "$size", Value: "$purchase_single_payments"}},
								bson.D{{Key: "$size", Value: "$purchase_monthly_payments"}},
							},
						},
					},
				},
			},
		}},
		// 4. Sort by purchase_count (descending)
		bson.D{{
			Key:   "$sort",
			Value: bson.D{{Key: "purchase_count", Value: -1}},
		}},
		// 5. Limit to top 10
		bson.D{{
			Key:   "$limit",
			Value: 10,
		}},
	}

	// Aggregate the top 10 cards by total number of purchases
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating top 10 cards: %w", err)
	}
	defer cursor.Close(context.TODO())

	var cards []models.Card
	for cursor.Next(context.TODO()) {
		var cardDoc entities.CardEntityNonSQL
		if err := cursor.Decode(&cardDoc); err != nil {
			return nil, fmt.Errorf("error decoding card doc: %w", err)
		}
		// Convert entity to model
		cards = append(cards, *entities.ToCard(&cardDoc))

	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return &cards, nil
}
