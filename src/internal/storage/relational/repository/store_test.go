package relational_repository

import (
	"log"
	"testing"

	mysql "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreWithHighestRevenueByMonth(t *testing.T) {
	logger.InitLogger(false, "")

	dsn := "testuser:testpassword@tcp(127.0.0.1:3306)/payment-registration-db?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := mysql.NewMySQLDB(dsn, true)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer mysql.CloseDB(database)

	err = mysql.ExecuteSQLFile(database, "../insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	storeRepo := NewStoreRelationalRepository(database)

	result, err := storeRepo.GetStoreWithHighestRevenueByMonth(10, 2024)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	assert.Equal(t, result.Cuit, "30-15066778-9")
	assert.Equal(t, result.Name, "Store O")
}
