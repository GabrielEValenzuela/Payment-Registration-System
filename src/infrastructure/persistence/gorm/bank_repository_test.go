package gorm

import (
	"testing"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/infrastructure/persistence/gorm/entities"
)

// Test the BankRepository with a real MySQL connection
func TestBankRepositoryWithMySQL(t *testing.T) {
	// Use the MySQL connection from mysql.go
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the models (create the tables if they don't exist)
	database.AutoMigrate(&entities.BankEntity{}, entities.FinancingEntity{})

	// Create a new BankRepository instance
	//_ = NewBankRepository(database)
	/*bankRepo := NewBankRepository(database)

	// Test data
	bank := models.BankModel{
		Cuit:      "30-12345678-9",
		Address:   "123 Main St",
		Telephone: "+54 11 1234 5678",
	}

	// Add the bank to the repository
	err = bankRepo.AddBank(bank)
	assert.NoError(t, err, "Error adding bank")

	// Fetch the bank from the repository
	var fetchedBank models.BankModel
	err = database.First(&fetchedBank, "cuit = ?", "30-12345678-9").Error
	assert.NoError(t, err, "Error fetching bank from database")

	// Assert that the bank was correctly inserted
	assert.Equal(t, bank.Cuit, fetchedBank.Cuit)
	assert.Equal(t, bank.Address, fetchedBank.Address)
	assert.Equal(t, bank.Telephone, fetchedBank.Telephone)*/
}
