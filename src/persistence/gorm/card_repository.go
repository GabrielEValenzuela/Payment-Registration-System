package gorm

import (
	"fmt"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"gorm.io/gorm"
)

type CardRepositoryGORM struct {
	db *gorm.DB
}

// NewCardRepository crea una nueva instancia de CardRepository
func NewCardRepository(db *gorm.DB) *CardRepositoryGORM {
	return &CardRepositoryGORM{db: db}
}

func (r *CardRepositoryGORM) GetPaymentSummary(cardNumber string, month int, year int) error {
	// Retrieve a card and its purchases in a specific month
	var card entities.CardEntity

	// Define the date range for the given month and year
	// startDate is the first day of the month, and endDate is the first day of the following month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // One month later

	// Query the card by its number and load purchases that match the month and year
	if err := r.db.Where("number = ?", cardNumber).
		Preload("PurchaseSinglePayments", "created_at >= ? AND created_at < ?", startDate, endDate).
		Preload("PurchaseMonthlyPayments", "created_at >= ? AND created_at < ?", startDate, endDate).
		First(&card).Error; err != nil {
		return nil
	}

	// Calculate the total purchases in that month
	var totalPrice float64
	for _, purchase := range card.PurchaseSinglePayments {
		totalPrice += purchase.PurchaseEntity.Amount
	}
	for _, purchase := range card.PurchaseMonthlyPayments {
		totalPrice += purchase.PurchaseEntity.Amount
	}

	// Define expiration dates
	firstExpiration := time.Now().AddDate(0, 0, 15)       // 15 days from today
	secondExpiration := firstExpiration.AddDate(0, 0, 10) // 10 days from today

	// Generate a unique code for the Payment Summary
	code := fmt.Sprintf("SUMMARY-%d-%d", year, month)

	// Create the PaymentSummaryEntity object
	paymentSummary := entities.PaymentSummaryEntity{
		Code:                code,
		Month:               month,
		Year:                year,
		FirstExpiration:     firstExpiration,
		SecondExpiration:    secondExpiration,
		SurchargePercentage: 5.0,        // Example: a 5% surcharge
		TotalPrice:          totalPrice, // Total of all purchases
		CardID:              card.ID,    // The card ID
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := r.db.Create(&paymentSummary).Error; err != nil {
		return fmt.Errorf("error inserting payment summary: %v", err)
	}

	return nil
}
