package relational_repository

import (
	"log"
	"testing"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	mysql "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailablePromotionsByStoreAndDateRange(t *testing.T) {
	testutils.InitTestSetup()

	// Use the MySQL connection from mysql.go
	dsn := testutils.DSN
	database, err := mysql.NewMySQLDB(dsn, true)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer mysql.CloseDB(database)

	// Insert Data
	err = mysql.ExecuteSQLFile(database, "../insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	// Test Financing Promotion
	testStore := "20-98765432-1"
	startDate := time.Date(2024, time.Month(10), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	promotionRepo := NewPromotionRelationRepository(database)

	financingPromotions, discountPromotions, err := promotionRepo.GetAvailablePromotionsByStoreAndDateRange(testStore, startDate, endDate)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(*discountPromotions))
	assert.Equal(t, 1, len(*financingPromotions))
}

func TestGetMostUsedPromotion(t *testing.T) {
	testutils.InitTestSetup()

	// Use the MySQL connection from mysql.go
	dsn := testutils.DSN
	database, err := mysql.NewMySQLDB(dsn, true)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer mysql.CloseDB(database)

	// Insert Data
	err = mysql.ExecuteSQLFile(database, "../insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	// Test Financing Promotion
	promotionRepo := NewPromotionRelationRepository(database)

	mostUsed, err := promotionRepo.GetMostUsedPromotion()
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}
	switch p := mostUsed.(type) {
	case entities.DiscountEntitySQL:
		log.Fatalf("Error")
	case entities.FinancingEntitySQL:
		assert.Equal(t, p.PromotionEntitySQL.Code, "PV20241001")
	default:
		log.Fatalf("Error")
	}
}
