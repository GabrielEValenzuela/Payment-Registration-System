package gorm

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/bank"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/mapper"
	testresource "github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/test_resource"
	"github.com/stretchr/testify/assert"
)

// Test the BankRepository with a real MySQL connection
func TestAddFinancingPromotionToBank(t *testing.T) {

	// Use the MySQL connection from mysql.go
	database, err := NewMySQLDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer CloseDB(database)

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

func TestExtendPromotionValidity(t *testing.T) {
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
	testCode := "PROMO123"
	testTime := time.Now().AddDate(0, 1, 1)

	bankRepo := NewBankRepository(database)

	bankRepo.ExtendFinancingPromotionValidity(testCode, testTime)

	var promotion entities.FinancingEntity
	if err := database.Where("code = ?", testCode).First(&promotion).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotion.Code, testCode)
	assert.Equal(t, promotion.ValidityEndDate.Unix(), testTime.Unix())

	// Test Discount Promotion
	testDiscountCode := "WINTERSALE2024"
	testDiscountTime := time.Now().AddDate(0, 1, 3)

	bankRepo.ExtendDiscountPromotionValidity(testDiscountCode, testDiscountTime)

	var promotionDiscount entities.DiscountEntity
	if err := database.Where("code = ?", testDiscountCode).First(&promotionDiscount).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotionDiscount.Code, testDiscountCode)
	assert.Equal(t, promotionDiscount.ValidityEndDate.Unix(), testDiscountTime.Unix())
}
