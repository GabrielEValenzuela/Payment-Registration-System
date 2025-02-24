package relational_repository

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	entities "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	mysql "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/testutils"
	"github.com/stretchr/testify/assert"
)

// Test the BankRepository with a real MySQL connection
func TestAddFinancingPromotionToBank(t *testing.T) {
	testutils.InitTestSetup()

	// Use the MySQL connection from mysql.go
	dsn := testutils.DSN
	database, err := mysql.NewMySQLDB(dsn, true)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer mysql.CloseDB(database)

	// insert new bank
	newBank := models.Bank{
		Name:      "Santander",
		Cuit:      "30-12345678-9",
		Address:   "123 Main St, Buenos Aires",
		Telephone: "+54 11 1234 5678",
	}

	result := database.Create(entities.ToBankEntity(&newBank))
	if result.Error != nil {
		fmt.Println("Failed to insert record:", result.Error)
	} else {
		fmt.Println("Bank inserted successfully")
	}

	// Insert Promotion
	newFinancingPromotion := models.Financing{
		Promotion: models.Promotion{
			Code:              "PROMO123",
			PromotionTitle:    "Summer Sale 2024",
			NameStore:         "Tech Store",
			CuitStore:         "30-12345678-9",
			ValidityStartDate: models.CustomTime{Time: time.Now().AddDate(0, -1, 0)}, // Mouth before
			ValidityEndDate:   models.CustomTime{Time: time.Now().AddDate(0, 1, 0)},  // Mouth after
			Comments:          "Special financing for summer purchases",
			Bank:              newBank,
		},
		NumberOfQuotas: 12,
		Interest:       5.5, // Tasa de interés
	}

	// Create a new BankRepository instance
	bankRepo := NewBankRelationalRepository(database)

	// Add the bank to the repository
	err = bankRepo.AddFinancingPromotionToBank(newFinancingPromotion)

	if err != nil {
		fmt.Println("Failed to insert record:", result.Error)
	} else {
		fmt.Println("Promotion inserted successfully")
	}
	// Fetch the bank from the repository
	var financingEntity entities.FinancingEntitySQL
	err = database.Preload("Bank").First(&financingEntity, "code = ?", "PROMO123").Error
	assert.NoError(t, err, "Error fetching promotion from database")

	// Assert that the bank was correctly inserted
	assert.Equal(t, newBank.Name, financingEntity.Bank.Name)
	assert.Equal(t, newFinancingPromotion.Code, financingEntity.Code)
	assert.Equal(t, newFinancingPromotion.Interest, financingEntity.Interest)
}

func TestExtendPromotionValidity(t *testing.T) {
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
	testCode := "PROMO123"
	testTime := time.Now().AddDate(0, 1, 1)

	bankRepo := NewBankRelationalRepository(database)

	bankRepo.ExtendFinancingPromotionValidity(testCode, testTime)

	var promotion entities.FinancingEntitySQL
	if err := database.Where("code = ?", testCode).First(&promotion).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotion.Code, testCode)
	assert.Equal(t, promotion.ValidityEndDate.Unix(), testTime.Unix())

	// Test Discount Promotion
	testDiscountCode := "WINTERSALE2024"
	testDiscountTime := time.Now().AddDate(0, 1, 3)

	bankRepo.ExtendDiscountPromotionValidity(testDiscountCode, testDiscountTime)

	var promotionDiscount entities.DiscountEntitySQL
	if err := database.Where("code = ?", testDiscountCode).First(&promotionDiscount).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotionDiscount.Code, testDiscountCode)
	assert.Equal(t, promotionDiscount.ValidityEndDate.Unix(), testDiscountTime.Unix())
}

func TestDeleteFinancingPromotion(t *testing.T) {
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
	testCode := "PV20241001"

	bankRepo := NewBankRelationalRepository(database)

	bankRepo.DeleteFinancingPromotion(testCode)

	var promotion entities.FinancingEntitySQL
	if err := database.Where("code = ?", testCode).First(&promotion).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotion.IsDeleted, true)
}

func TestDeleteDiscountPromotion(t *testing.T) {
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
	testCode := "SPRINGDEAL2024"

	bankRepo := NewBankRelationalRepository(database)

	bankRepo.DeleteDiscountPromotion(testCode)

	var promotion entities.DiscountEntitySQL
	if err := database.Where("code = ?", testCode).First(&promotion).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotion.IsDeleted, true)
}

func TestGetBankCustomerCounts(t *testing.T) {
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

	bankRepo := NewBankRelationalRepository(database)

	result, err := bankRepo.GetBankCustomerCounts()
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	assert.Equal(t, len(result), 4)

	var bank *models.BankCustomerCountDTO
	for _, v := range result {
		if v.BankName == "Santander" {
			bank = &v
		}
	}
	assert.Equal(t, bank.BankName, "Santander")
	assert.Equal(t, bank.CustomerCount, 2)
}
