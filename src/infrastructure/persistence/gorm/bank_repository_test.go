package gorm

import (
	"fmt"
	"testing"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/infrastructure/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/infrastructure/persistence/gorm/mapper"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/bank"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
	"github.com/stretchr/testify/assert"
)

// Test the BankRepository with a real MySQL connection
func TestBankRepositoryWithMySQL(t *testing.T) {
	// Use the MySQL connection from mysql.go
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the models (create the tables if they don't exist)
	database.Migrator().DropTable(&entities.BankEntity{}, entities.FinancingEntity{})
	database.AutoMigrate(&entities.BankEntity{}, entities.FinancingEntity{})

	// insert new bank
	newBank := bank.Bank{
		Name:      "Santander",
		Cuit:      "30-12345678-9",
		Address:   "123 Main St, Buenos Aires",
		Telephone: "+54 11 1234 5678",
	}

	result := database.Create(mapper.ToBankEntity(&newBank))
	if result.Error != nil {
		fmt.Println("Failed to insert record:", result.Error)
	} else {
		fmt.Println("Bank inserted successfully")
	}

	// Insert Promotion
	newFinancingPromotion := promotion.Financing{
		Promotion: promotion.Promotion{
			Code:              "PROMO123",
			PromotionTitle:    "Summer Sale 2024",
			NameStore:         "Tech Store",
			CuitStore:         "30-12345678-9",
			ValidityStartDate: time.Now().AddDate(0, -1, 0), // Un mes antes
			ValidityEndDate:   time.Now().AddDate(0, 1, 0),  // Un mes después
			Comments:          "Special financing for summer purchases",
			Bank:              newBank,
		},
		NumberOfQuotas: 12,
		Interest:       5.5, // Tasa de interés
	}

	// Create a new BankRepository instance
	bankRepo := NewBankRepository(database)

	// Add the bank to the repository
	err = bankRepo.AddFinancingPromotionToBank(&newFinancingPromotion)

	if err != nil {
		fmt.Println("Failed to insert record:", result.Error)
	} else {
		fmt.Println("Promotion inserted successfully")
	}
	// Fetch the bank from the repository
	var financingEntity entities.FinancingEntity
	err = database.Preload("Bank").First(&financingEntity, "code = ?", "PROMO123").Error
	assert.NoError(t, err, "Error fetching promotion from database")

	// Assert that the bank was correctly inserted
	assert.Equal(t, newBank.Name, financingEntity.Bank.Name)
	assert.Equal(t, newFinancingPromotion.Code, financingEntity.Code)
	assert.Equal(t, newFinancingPromotion.Interest, financingEntity.Interest)
}
