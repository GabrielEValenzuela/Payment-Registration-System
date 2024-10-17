package gorm

import (
	"log"
	"testing"
	"time"

	testresource "github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/test_resource"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailablePromotionsByStoreAndDateRange(t *testing.T) {
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

	// Test Financing Promotion
	testStore := "20-98765432-1"
	startDate := time.Date(2024, time.Month(10), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	promotionRepo := NewPromotionRepository(database)

	financingPromotions, discountPromotions, err := promotionRepo.GetAvailablePromotionsByStoreAndDateRange(testStore, startDate, endDate)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(*discountPromotions))
	assert.Equal(t, 1, len(*financingPromotions))
}
