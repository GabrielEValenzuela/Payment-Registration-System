package gorm

import (
	"fmt"
	"log"
	"testing"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	testresource "github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/test_resource"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePaymentSummary(t *testing.T) {
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

	cardRepo.GetPaymentSummary(cardNumber, month, year)

	// Assert
	code := fmt.Sprintf("SUMMARY-%d-%d", year, month)

	var paymentSummary entities.PaymentSummaryEntity
	if err := database.Where("code = ?", code).First(&paymentSummary).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", code, err))
	}

	assert.Equal(t, paymentSummary.Code, code)
	assert.Equal(t, paymentSummary.TotalPrice, 210.00)
}
