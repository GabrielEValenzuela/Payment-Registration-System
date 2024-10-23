package gorm

import (
	"log"
	"testing"

	testresource "github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/test_resource"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreWithHighestRevenueByMonth(t *testing.T) {
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

	err = testresource.ExecuteSQLFile(database, "./test_resource/insert.sql")
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	storeRepo := NewStoreRepository(database)

	result, err := storeRepo.GetStoreWithHighestRevenueByMonth(10, 2024)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	assert.Equal(t, result.Cuit, "30-15066778-9")
	assert.Equal(t, result.Name, "Store O")
}
