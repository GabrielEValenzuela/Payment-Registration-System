package gorm

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/purchase"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/mapper"
	testresource "github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/test_resource"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentSummary(t *testing.T) {
	cardNumber := "1234567812345678"
	month := 10
	year := 2024

	// Database
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

	// Insert Data
	err = testresource.ExecuteSQLFile(database, "./test_resource/insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	// Test Repository
	cardRepo := NewCardRepository(database)

	paymentSummary, err := cardRepo.GetPaymentSummary(cardNumber, month, year)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	// Assert
	code := fmt.Sprintf("SUMMARY-%d-%d", year, month)

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // One month later

	var paymentSummaryEntity entities.PaymentSummaryEntity
	if err := database.
		Joins("JOIN CARDS ON CARDS.id = PAYMENT_SUMMARIES.card_id").
		Preload("Card").
		Preload("Card.PurchaseSinglePayments", "created_at >= ? AND created_at < ?", startDate, endDate).
		Preload("Card.PurchaseMonthlyPayments", "created_at >= ? AND created_at < ?", startDate, endDate).
		Where("number = ?", cardNumber).
		Where("code = ?", code).First(&paymentSummaryEntity).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", code, err))
	}

	assert.Equal(t, paymentSummaryEntity.Code, code)
	assert.Equal(t, paymentSummaryEntity.TotalPrice, 510.00)
	assert.Equal(t, len(paymentSummaryEntity.Card.PurchaseMonthlyPayments), 1)
	assert.Equal(t, len(paymentSummaryEntity.Card.PurchaseSinglePayments), 3)

	paymentSummaryMapper := mapper.ToPaymentSummary(&paymentSummaryEntity)

	assert.Equal(t, paymentSummaryMapper.Code, paymentSummary.Code)
	assert.Equal(t, paymentSummaryMapper.TotalPrice, paymentSummary.TotalPrice)
	assert.Equal(t, len(paymentSummaryMapper.Card.PurchaseMonthlyPayments), len(paymentSummary.Card.PurchaseMonthlyPayments))
	assert.Equal(t, len(paymentSummaryMapper.Card.PurchaseSinglePayments), len(paymentSummary.Card.PurchaseSinglePayments))
}

func TestGetCardsExpiringInNext30Days(t *testing.T) {
	day := 16
	month := 10
	year := 2024

	// Database
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

	// Insert Data
	err = testresource.ExecuteSQLFile(database, "./test_resource/insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	cardRepo := NewCardRepository(database)

	cards, err := cardRepo.GetCardsExpiringInNext30Days(day, month, year)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	assert.Equal(t, len(*cards), 4)
}

func TestGetPurchaseSingle(t *testing.T) {
	paymentVoucher := "PV20241001"
	cuit := "30-12345678-9"
	finalAmount := 100.00

	// Database
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

	// Insert Data
	err = testresource.ExecuteSQLFile(database, "./test_resource/insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	cardRepo := NewCardRepository(database)

	payment, err := cardRepo.GetPurchaseSingle(cuit, finalAmount, paymentVoucher)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	assert.Equal(t, payment.Purchase.Store, "Store A")
}

func TestGetPurchaseMonthly(t *testing.T) {
	paymentVoucher := "PV20241101"
	cuit := "20-98765432-1"
	finalAmount := 440.0

	// Database
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

	// Insert Data
	err = testresource.ExecuteSQLFile(database, "./test_resource/insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	cardRepo := NewCardRepository(database)

	payment, err := cardRepo.GetPurchaseMonthly(cuit, finalAmount, paymentVoucher)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	assert.Equal(t, payment.Purchase.Store, "Store B")
	assert.Equal(t, payment.Purchase.Amount, 110.00)
	assert.Equal(t, len(payment.Quota), 4)
}

func TestGetTop10CardsByPurchases(t *testing.T) {
	// Database
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

	// Insert Data
	err = testresource.ExecuteSQLFile(database, "./test_resource/insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	cardRepo := NewCardRepository(database)

	cards, err := cardRepo.GetTop10CardsByPurchases()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	assert.Equal(t, len(*cards), 10)

	var card *card.Card
	for _, src := range *cards {
		if src.Number == "1234567812345678" {
			card = &src
		}
	}

	assert.Equal(t, len(card.PurchaseSinglePayments), 5)
	assert.Equal(t, len(card.PurchaseMonthlyPayments), 2)

	var purchaseMonthly *purchase.PurchaseMonthlyPayment
	for _, src := range card.PurchaseMonthlyPayments {
		if src.Purchase.PaymentVoucher == "PV20241001" {
			purchaseMonthly = &src
		}
	}
	assert.Equal(t, len(purchaseMonthly.Quota), 3)
}
