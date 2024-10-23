package gorm

import (
	"fmt"
	"log"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/payment_summary"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/purchase"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/mapper"
	"gorm.io/gorm"
)

type CardRepositoryGORM struct {
	db *gorm.DB
}

// NewCardRepository crea una nueva instancia de CardRepository
func NewCardRepository(db *gorm.DB) *CardRepositoryGORM {
	return &CardRepositoryGORM{db: db}
}

func (r *CardRepositoryGORM) GetPaymentSummary(cardNumber string, month int, year int) (*payment_summary.PaymentSummary, error) {
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
		return nil, err
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
		return nil, fmt.Errorf("error inserting payment summary: %v", err)
	}

	paymentSummary.Card = card

	return mapper.ToPaymentSummary(&paymentSummary), nil
}

func (r *CardRepositoryGORM) GetCardsExpiringInNext30Days(day int, month int, year int) (*[]card.Card, error) {
	startDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	next30Days := startDate.AddDate(0, 0, 30)

	var paymentSummaryList []entities.PaymentSummaryEntity

	if err := r.db.Preload("Card").Where("first_expiration BETWEEN ? AND ?", startDate, next30Days).Find(&paymentSummaryList).Error; err != nil {
		return nil, err
	}

	var cards []card.Card
	for _, src := range paymentSummaryList {
		cards = append(cards, *mapper.ToCard(&src.Card))
	}

	return &cards, nil
}

func (r *CardRepositoryGORM) GetPurchaseSingle(cuit string, finalAmount float64, paymentVoucher string) (*purchase.PurchaseSinglePayment, error) {
	var paymentEntity entities.PurchaseSinglePaymentEntity

	if err := r.db.Where("cuit_store = ? AND final_amount = ? AND payment_voucher = ?", cuit, finalAmount, paymentVoucher).
		First(&paymentEntity).Error; err != nil {
		return nil, err
	}

	return mapper.ToPurchaseSinglePayment(&paymentEntity), nil
}

func (r *CardRepositoryGORM) GetPurchaseMonthly(cuit string, finalAmount float64, paymentVoucher string) (*purchase.PurchaseMonthlyPayment, error) {
	var paymentEntity entities.PurchaseMonthlyPaymentsEntity

	if err := r.db.Preload("Quotas").Where("cuit_store = ? AND final_amount = ? AND payment_voucher = ?", cuit, finalAmount, paymentVoucher).
		First(&paymentEntity).Error; err != nil {
		return nil, err
	}

	return mapper.ToPurchaseMonthlyPayments(&paymentEntity), nil
}

func (r *CardRepositoryGORM) GetTop10CardsByPurchases() (*[]card.Card, error) {
	var cardEntities []entities.CardEntity

	// Query to find the top 10 cards by number of purchases, with preloads for payments and quotas
	if err := r.db.Table("CARDS").
		Select("CARDS.*, (COUNT(PURCHASES_SINGLE_PAYMENTS.id) + COUNT(PURCHASES_MONTHLY_PAYMENTS.id)) as purchase_count").
		Joins("LEFT JOIN PURCHASES_SINGLE_PAYMENTS ON PURCHASES_SINGLE_PAYMENTS.card_id = CARDS.id").
		Joins("LEFT JOIN PURCHASES_MONTHLY_PAYMENTS ON PURCHASES_MONTHLY_PAYMENTS.card_id = CARDS.id").
		Preload("PurchaseSinglePayments").
		Preload("PurchaseMonthlyPayments.Quotas").
		Group("CARDS.id").
		Order("purchase_count DESC").
		Limit(10).
		Find(&cardEntities).Error; err != nil {
		log.Printf("Error retrieving top 10 cards by purchases: %v", err)
		return nil, err
	}

	var cards []card.Card

	for _, card := range cardEntities {
		cards = append(cards, *mapper.ToCard(&card))
	}

	return &cards, nil
}
