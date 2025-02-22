/*
 * Payment Registration System - Component Tests
 * ---------------------------------------------
 * This file contains the component tests for the payment registration system.
 * It includes tests for the CRUD for all endpoints, as well as the integration tests for the system.
 *
 * Created: Feb. 22, 2025
 * License: GNU General Public License v3.0
 */
package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	nonrelational "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/non_relational"
	non_relational_repository "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/non_relational/repository"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational"
	relational_repository "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational/repository"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

// Aux function to get the absolute path of a file
func getAbsolutePath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0) // Get the current file path
	basePath := filepath.Dir(filename)     // Get the directory of the current file
	absolutePath, err := filepath.Abs(filepath.Join(basePath, relativePath))
	if err != nil {
		panic("Failed to determine absolute path: " + err.Error())
	}
	return absolutePath
}

var (
	SQLDatabase   *gorm.DB        // MySQL (GORM) database instance
	NoSQLDatabase *mongo.Database // MongoDB database instance
	// Define paths dynamically
	relationalSQLPath = getAbsolutePath("data/relational/database.sql")
)

const (
	mysqlDSN      = "root:root@tcp(127.0.0.1:3306)/"
	mysqlDatabase = "payment_registration_system"
	mysqlUser     = "testuser"
	mysqlPassword = "testpassword"
	mongoURI      = "mongodb://root:root@localhost:27017/"
	mongoDatabase = "payment_registration_system"
	URI_MONGO     = "mongodb://testuser:testpassword@localhost:27017/payment_registration_system?authSource=admin"
	MYSQL_DSN     = "testuser:testpassword@tcp(127.0.0.1:3306)/payment_registration_system?charset=utf8mb4&parseTime=True&loc=Local"
)

// Aux function to load the SQL data from a file
func executeSQLFile(db *gorm.DB, filePath string) error {
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %w", err)
	}

	fmt.Printf("Executing SQL file %s\n", filePath)

	queries := strings.Split(string(sqlContent), ";\n")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" {
			if err := db.Exec(query).Error; err != nil {
				return fmt.Errorf("error executing query: %v", err)
			}
		}
	}

	return nil
}

// Aux function to initialize the test setup
func InitTestSetup() {
	var once sync.Once
	once.Do(func() {
		logger.InitLogger(false, "test.log")
	})

	var err error
	SQLDatabase, err = relational.NewMySQLDB(MYSQL_DSN, true)
	if err != nil {
		log.Fatalf("❌ Error initializing MySQL database: %v", err)
	}

	NoSQLDatabase, err = nonrelational.NewMongoDB(URI_MONGO, mongoDatabase, false)
	if err != nil {
		log.Fatalf("❌ Error initializing MongoDB database: %v", err)
	}

	if err := executeSQLFile(SQLDatabase, relationalSQLPath); err != nil {
		log.Fatalf("❌ Error importing SQL data: %v", err)
	}
}

func TestMain(m *testing.M) {
	// Initialize databases once before running all tests
	InitTestSetup()

	// Run all tests
	exitCode := m.Run()

	// Close MySQL connection
	err := relational.CloseDB(SQLDatabase)
	if err != nil {
		log.Fatalf("Error closing MySQL database: %v", err)
	}
	// CleanMongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info("Cleaning the database: dropping existing collections...")
	collections, err := NoSQLDatabase.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		fmt.Printf("failed to list collections: %v", err)
	}

	for _, collection := range collections {
		if err := NoSQLDatabase.Collection(collection).Drop(ctx); err != nil {
			fmt.Printf("failed to drop collection %s: %v", collection, err)
		}
		logger.Info("Dropped collection: %s", collection)
	}

	// Close MongoDB connection
	err = nonrelational.CloseMongoDB(NoSQLDatabase.Client())
	if err != nil {
		log.Fatalf("Error closing MongoDB: %v", err)
	}

	// Exit
	os.Exit(exitCode)
}

// Tests

// ---------------------------------------------------
// ---------------     BANK TESTS --------------------
// ---------------------------------------------------
func TestAddFinancingPromotionToBank(t *testing.T) {
	var codeString = "PROMO123TEST"
	// Create a new bank
	newBank := models.Bank{
		Name:      "Santander",
		Cuit:      "30-12345678-9",
		Address:   "123 Main St, Buenos Aires",
		Telephone: "+54 11 1234 5678",
	}

	// Insert the new bank
	result := SQLDatabase.Create(entities.ToBankEntity(&newBank))
	assert.NoError(t, result.Error, "Error inserting bank into MySQL")

	// Insert a promotion
	newFinancingPromotion := models.Financing{
		Promotion: models.Promotion{
			Code:              codeString,
			PromotionTitle:    "Summer Sale 2024",
			NameStore:         "Tech Store",
			CuitStore:         "30-12345678-9",
			ValidityStartDate: models.CustomTime{Time: time.Now().AddDate(0, -1, 0)}, // 1 month before
			ValidityEndDate:   models.CustomTime{Time: time.Now().AddDate(0, 1, 0)},  // 1 month after
			Comments:          "Special financing for summer purchases",
			Bank:              newBank,
		},
		NumberOfQuotas: 12,
		Interest:       5.5, // Interest rate
	}

	// ------ SQL (MySQL) ------

	// Use BankRelationalRepository to insert the financing promotion
	bankRepo := relational_repository.NewBankRelationalRepository(SQLDatabase)
	err := bankRepo.AddFinancingPromotionToBank(newFinancingPromotion)
	assert.NoError(t, err, "Error inserting financing promotion into MySQL")

	// Verify that the promotion was inserted
	var financingEntity entities.FinancingEntitySQL
	err = SQLDatabase.Preload("Bank").First(&financingEntity, "code = ?", codeString).Error
	assert.NoError(t, err, "Error fetching promotion from MySQL")

	// Assertions
	assert.Equal(t, newBank.Name, financingEntity.Bank.Name)
	assert.Equal(t, newFinancingPromotion.Code, financingEntity.Code)
	assert.Equal(t, newFinancingPromotion.Interest, financingEntity.Interest)

	// ------ NoSQL (MongoDB) ------

	// Use MongoDB repository
	noSQLBankRepo := non_relational_repository.NewBankNonRelationalRepository(NoSQLDatabase)
	err = noSQLBankRepo.AddFinancingPromotionToBank(newFinancingPromotion)
	assert.NoError(t, err, "Error inserting financing promotion into MongoDB")

	var resultMongo struct {
		Financing entities.FinancingEntityNonSQL `bson:",inline"`
		Bank      entities.BankEntityNonSQL      `bson:"bank"`
	}

	pipeline := []bson.M{
		{"$match": bson.M{"promotion_entity.code": codeString}},
		{
			"$lookup": bson.M{
				"from":         "banks",
				"localField":   "bank_id",
				"foreignField": "_id",
				"as":           "bank",
			},
		},
		{"$unwind": "$bank"}, // Ensure only one bank object is returned
	}

	cursor, err := NoSQLDatabase.Collection("financings").Aggregate(context.TODO(), pipeline)
	assert.NoError(t, err, "Error fetching promotion from MongoDB")
	assert.True(t, cursor.Next(context.TODO()), "No financing promotion found")

	// Decode query result into resultMongo
	err = cursor.Decode(&resultMongo)
	assert.NoError(t, err, "Error decoding promotion data")

	// Assertions
	assert.Equal(t, newBank.Name, resultMongo.Bank.Name, "Bank name mismatch in MongoDB")
	assert.Equal(t, newBank.Cuit, resultMongo.Bank.Cuit, "Bank CUIT mismatch in MongoDB")
	assert.Equal(t, newFinancingPromotion.Promotion.Code, resultMongo.Financing.PromotionEntity.Code, "Code mismatch in MongoDB")
	assert.Equal(t, newFinancingPromotion.Interest, resultMongo.Financing.Interest, "Interest mismatch in MongoDB")

}

func TestBankExtendPromotionValidity(t *testing.T) {
	// TestBank Financing Promotion
	testCode := "PROMO123TEST"
	testTime := time.Now().AddDate(0, 1, 1)

	bankRepo := relational_repository.NewBankRelationalRepository(SQLDatabase)

	err := bankRepo.ExtendFinancingPromotionValidity(testCode, testTime)
	assert.NoError(t, err, "Error extending financing promotion validity in MySQL")

	// Validate in SQL Database
	var promotion entities.FinancingEntitySQL
	err = SQLDatabase.Where("code = ?", testCode).First(&promotion).Error
	assert.NoError(t, err, fmt.Sprintf("Could not find financing promotion with code %s", testCode))
	assert.Equal(t, testCode, promotion.Code)
	assert.Equal(t, testTime.Unix(), promotion.ValidityEndDate.Unix())

	// TestBank Discount Promotion
	testDiscountCode := "WINTERSALE2024"
	testDiscountTime := time.Now().AddDate(0, 1, 3)

	err = bankRepo.ExtendDiscountPromotionValidity(testDiscountCode, testDiscountTime)
	assert.NoError(t, err, "Error extending discount promotion validity in MySQL")

	// Validate in SQL Database
	var promotionDiscount entities.DiscountEntitySQL
	err = SQLDatabase.Where("code = ?", testDiscountCode).First(&promotionDiscount).Error
	assert.NoError(t, err, fmt.Sprintf("Could not find discount promotion with code %s", testDiscountCode))
	assert.Equal(t, testDiscountCode, promotionDiscount.Code)
	assert.Equal(t, testDiscountTime.Unix(), promotionDiscount.ValidityEndDate.Unix())

	// ------ NoSQL (MongoDB) ------

	// Use MongoDB repository
	noSQLBankRepo := non_relational_repository.NewBankNonRelationalRepository(NoSQLDatabase)

	err = noSQLBankRepo.ExtendFinancingPromotionValidity(testCode, testTime)
	fmt.Println(err)
	assert.NoError(t, err, "Error extending financing promotion validity in MongoDB")

	// Validate in MongoDB
	var resultMongo struct {
		Financing entities.FinancingEntityNonSQL `bson:",inline"`
	}
	pipeline := []bson.M{
		{"$match": bson.M{"promotion_entity.code": testCode}},
	}

	cursor, err := NoSQLDatabase.Collection("financings").Aggregate(context.TODO(), pipeline)
	assert.NoError(t, err, "Error fetching financing promotion from MongoDB")
	assert.True(t, cursor.Next(context.TODO()), "No financing promotion found in MongoDB")
	err = cursor.Decode(&resultMongo)
	assert.NoError(t, err, "Error decoding financing promotion from MongoDB")
}

func TestBankDeleteFinancingPromotion(t *testing.T) {
	// TestBank Financing Promotion
	testCode := "PV20241001"

	bankRepo := relational_repository.NewBankRelationalRepository(SQLDatabase)

	bankRepo.DeleteFinancingPromotion(testCode)

	var promotion entities.FinancingEntitySQL
	if err := SQLDatabase.Where("code = ?", testCode).First(&promotion).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotion.IsDeleted, true)

	// ------ NoSQL (MongoDB) ------

	// Use MongoDB repository
	noSQLBankRepo := non_relational_repository.NewBankNonRelationalRepository(NoSQLDatabase)

	err := noSQLBankRepo.DeleteFinancingPromotion(testCode)
	assert.NoError(t, err, "Error deleting financing promotion in MongoDB")
}

func TestBankDeleteDiscountPromotion(t *testing.T) {
	// Test Financing Promotion
	testCode := "SPRINGDEAL2024"

	bankRepo := relational_repository.NewBankRelationalRepository(SQLDatabase)

	bankRepo.DeleteDiscountPromotion(testCode)

	var promotion entities.DiscountEntitySQL
	if err := SQLDatabase.Where("code = ?", testCode).First(&promotion).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", testCode, err))
	}

	assert.Equal(t, promotion.IsDeleted, true)

	// ------ NoSQL (MongoDB) ------

	// Use MongoDB repository
	noSQLBankRepo := non_relational_repository.NewBankNonRelationalRepository(NoSQLDatabase)
	err := noSQLBankRepo.DeleteDiscountPromotion(testCode)
	assert.NoError(t, err, "Error deleting discount promotion in MongoDB")

}

func TestBankGetBankCustomerCounts(t *testing.T) {
	bankRepo := relational_repository.NewBankRelationalRepository(SQLDatabase)

	result, err := bankRepo.GetBankCustomerCounts()
	assert.NoError(t, err, "Error fetching bank customer counts from MySQL")

	assert.Greater(t, len(result), 0)

	var bank *models.BankCustomerCountDTO
	for _, v := range result {
		if v.BankName == "Santander" {
			bank = &v
		}
	}
	assert.Equal(t, bank.BankName, "Santander")

	// ------ NoSQL (MongoDB) ------

	// Use MongoDB repository
	noSQLBankRepo := non_relational_repository.NewBankNonRelationalRepository(NoSQLDatabase)
	resultMongo, err := noSQLBankRepo.GetBankCustomerCounts()

	assert.NoError(t, err, "Error fetching bank customer counts from MongoDB")

	assert.Greater(t, len(resultMongo), 0)

	var bankMongo *models.BankCustomerCountDTO
	for _, v := range resultMongo {
		if v.BankName == "Santander" {
			bankMongo = &v
		}
	}

	assert.Equal(t, bankMongo.BankName, "Santander")
}

// ---------------------------------------------------
// ---------------     CARD TESTS --------------------
// ---------------------------------------------------
func TestCardGetPaymentSummary(t *testing.T) {
	cardNumber := "1234567812345678"
	month := 10
	year := 2024

	// ------ SQL (MySQL) ------
	cardRepo := relational_repository.NewCardRelationalRepository(SQLDatabase)

	paymentSummary, err := cardRepo.GetPaymentSummary(cardNumber, month, year)
	assert.NoError(t, err, "Error fetching payment summary from MySQL")

	// ✅ Generate expected code format
	expectedCode := fmt.Sprintf("SUMMARY-%d-%d", year, month)

	// ✅ Define start & end date for the query
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // One month later

	// ✅ Fetch from MySQL to validate correctness
	var paymentSummaryEntity entities.PaymentSummaryEntitySQL
	if err := SQLDatabase.
		Joins("JOIN CARDS ON CARDS.id = PAYMENT_SUMMARIES.card_id").
		Preload("Card").
		Preload("Card.PurchaseSinglePayments", "created_at >= ? AND created_at < ?", startDate, endDate).
		Preload("Card.PurchaseMonthlyPayments", "created_at >= ? AND created_at < ?", startDate, endDate).
		Where("number = ?", cardNumber).
		Where("code = ?", expectedCode).First(&paymentSummaryEntity).Error; err != nil {
		panic(fmt.Errorf("could not find promotion with code %s: %v", expectedCode, err))
	}

	assert.Equal(t, paymentSummaryEntity.Code, expectedCode)
	assert.Equal(t, paymentSummaryEntity.TotalPrice, 510.00)
	assert.Equal(t, len(paymentSummaryEntity.Card.PurchaseMonthlyPayments), 1)
	assert.Equal(t, len(paymentSummaryEntity.Card.PurchaseSinglePayments), 3)

	paymentSummaryMapper := entities.ToPaymentSummary(&paymentSummaryEntity)

	assert.Equal(t, paymentSummaryMapper.Code, paymentSummary.Code)
	assert.Equal(t, paymentSummaryMapper.TotalPrice, paymentSummary.TotalPrice)
	assert.Equal(t, len(paymentSummaryMapper.Card.PurchaseMonthlyPayments), len(paymentSummary.Card.PurchaseMonthlyPayments))
	assert.Equal(t, len(paymentSummaryMapper.Card.PurchaseSinglePayments), len(paymentSummary.Card.PurchaseSinglePayments))

	// ------ NoSQL (MongoDB) ------
	noSQLCardRepo := non_relational_repository.NewCardNonRelationalRepository(NoSQLDatabase)
	paymentSummaryMongo, err := noSQLCardRepo.GetPaymentSummary(cardNumber, month, year)

	assert.NoError(t, err, "Error fetching payment summary from MongoDB")

	assert.Equal(t, paymentSummaryMongo.Code, expectedCode)
	assert.Equal(t, 420.0, paymentSummaryMongo.TotalPrice)
	assert.Equal(t, 1, len(paymentSummaryMongo.SinglePayments))
	assert.Equal(t, 1, len(paymentSummaryMongo.MonthlyPayments))

}

func TestCardGetCardsExpiringInNext30Days(t *testing.T) {
	day := 16
	month := 10
	year := 2024

	cardRepo := relational_repository.NewCardRelationalRepository(SQLDatabase)

	cards, err := cardRepo.GetCardsExpiringInNext30Days(day, month, year)
	assert.NoError(t, err, "Error fetching cards expiring in the next 30 days from MySQL")

	assert.Equal(t, 4, len(*cards))

	// ------ NoSQL (MongoDB) ------
	noSQLCardRepo := non_relational_repository.NewCardNonRelationalRepository(NoSQLDatabase)
	cardsMongo, err := noSQLCardRepo.GetCardsExpiringInNext30Days(day, month, year)

	assert.NoError(t, err, "Error fetching cards expiring in the next 30 days from MongoDB")

	assert.Equal(t, 2, len(*cardsMongo))
}

func TestCardGetPurchaseSingle(t *testing.T) {
	paymentVoucher := "PV20241001"
	cuit := "30-12345678-9"
	finalAmount := 100.00

	cardRepo := relational_repository.NewCardRelationalRepository(SQLDatabase)

	payment, err := cardRepo.GetPurchaseSingle(cuit, finalAmount, paymentVoucher)
	assert.NoError(t, err, "Error fetching purchase single from MySQL")

	assert.Equal(t, payment.Purchase.Store, "Store A")

	// ------ NoSQL (MongoDB) ------
	finalAmount = 90
	noSQLCardRepo := non_relational_repository.NewCardNonRelationalRepository(NoSQLDatabase)
	paymentMongo, err := noSQLCardRepo.GetPurchaseSingle(cuit, finalAmount, paymentVoucher)

	assert.NoError(t, err, "Error fetching purchase single from MongoDB")

	assert.Equal(t, paymentMongo.Store, "Store A")
}

func TestCardGetPurchaseMonthly(t *testing.T) {
	paymentVoucher := "PV20241101"
	cuit := "20-98765432-1"
	finalAmount := 440.0

	cardRepo := relational_repository.NewCardRelationalRepository(SQLDatabase)

	payment, err := cardRepo.GetPurchaseMonthly(cuit, finalAmount, paymentVoucher)
	assert.NoError(t, err, "Error fetching purchase monthly from MySQL")

	assert.Equal(t, payment.Purchase.Store, "Store B")
	assert.Equal(t, payment.Purchase.Amount, 110.00)
	assert.Equal(t, len(payment.Quota), 4)

	// ------ NoSQL (MongoDB) ------
	paymentVoucher = "PV20241001"
	cuit = "30-12345678-9"
	finalAmount = 330.00

	noSQLCardRepo := non_relational_repository.NewCardNonRelationalRepository(NoSQLDatabase)
	paymentMongo, err := noSQLCardRepo.GetPurchaseMonthly(cuit, finalAmount, paymentVoucher)

	assert.NoError(t, err, "Error fetching purchase single from MongoDB")

	assert.Equal(t, paymentMongo.Store, "Store A")
	assert.Equal(t, paymentMongo.Amount, 300.00)
	assert.Equal(t, len(paymentMongo.Quota), 3)
	assert.Equal(t, paymentMongo.Quota[0].Price, 110.00)
}
func TestCardGetTop10CardsByPurchases(t *testing.T) {

	cardRepo := relational_repository.NewCardRelationalRepository(SQLDatabase)

	cards, err := cardRepo.GetTop10CardsByPurchases()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	assert.Equal(t, len(*cards), 10)

	var card *models.Card
	for _, src := range *cards {
		if src.Number == "1234567812345678" {
			card = &src
		}
	}

	assert.Equal(t, len(card.PurchaseSinglePayments), 5)
	assert.Equal(t, len(card.PurchaseMonthlyPayments), 2)

	var purchaseMonthly *models.PurchaseMonthlyPayment
	for _, src := range card.PurchaseMonthlyPayments {
		if src.Purchase.PaymentVoucher == "PV20241001" {
			purchaseMonthly = &src
		}
	}
	assert.Equal(t, len(purchaseMonthly.Quota), 3)

	// ------ NoSQL (MongoDB) ------
	noSQLCardRepo := non_relational_repository.NewCardNonRelationalRepository(NoSQLDatabase)
	cardsMongo, err := noSQLCardRepo.GetTop10CardsByPurchases()

	assert.NoError(t, err, "Error fetching top 10 cards by purchases from MongoDB")

	assert.Equal(t, 10, len(*cardsMongo))

	var cardMongo *models.Card
	for _, src := range *cardsMongo {
		if src.Number == "7446548631079191" {
			cardMongo = &src
			fmt.Printf("Found card: %s belong to %s\n", cardMongo.Number, cardMongo.CardholderNameInCard)
			break
		}
	}

	assert.Equal(t, len(cardMongo.PurchaseSinglePayments), 2)
	assert.Equal(t, len(cardMongo.PurchaseMonthlyPayments), 2)

	var purchaseMonthlyMongo *models.PurchaseMonthlyPayment
	for _, src := range cardMongo.PurchaseMonthlyPayments {
		if src.Purchase.PaymentVoucher == "PV-M581450" {
			fmt.Printf("Found purchase monthly: %s\n", src.Purchase.PaymentVoucher)
			purchaseMonthlyMongo = &src

		}
	}

	assert.Equal(t, 12, len(purchaseMonthlyMongo.Quota))
	assert.Equal(t, float64(4), purchaseMonthlyMongo.Interest)
	assert.Equal(t, 53.3, purchaseMonthlyMongo.Quota[0].Price)
	assert.Equal(t, int64(purchaseMonthlyMongo.FinalAmount), int64(float64(len(purchaseMonthlyMongo.Quota))*purchaseMonthlyMongo.Quota[0].Price))

}

// ---------------------------------------------------
// --------------  PROMOTION TESTS -------------------
// ---------------------------------------------------

func TestPromotionGetAvailablePromotionsByStoreAndDateRange(t *testing.T) {
	// Test Financing Promotion
	testStore := "30-12345678-9"
	startDate := time.Date(2024, time.Month(10), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	promotionRepo := relational_repository.NewPromotionRelationRepository(SQLDatabase)

	financingPromotions, discountPromotions, err := promotionRepo.GetAvailablePromotionsByStoreAndDateRange(testStore, startDate, endDate)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(*discountPromotions))
	assert.Equal(t, 1, len(*financingPromotions))

	// ------ NoSQL (MongoDB) ------
	noSQLPromotionRepo := non_relational_repository.NewPromotionNonRelationalRepository(NoSQLDatabase)
	financingPromotionsMongo, discountPromotionsMongo, err := noSQLPromotionRepo.GetAvailablePromotionsByStoreAndDateRange(testStore, startDate, endDate)

	assert.NoError(t, err, "Error fetching available promotions by store and date range from MongoDB")

	assert.Equal(t, 1, len(*discountPromotionsMongo))
	assert.Equal(t, 1, len(*financingPromotionsMongo))
}

func TestPromotionGetMostUsedPromotion(t *testing.T) {
	// Test Financing Promotion
	promotionRepo := relational_repository.NewPromotionRelationRepository(SQLDatabase)

	mostUsed, err := promotionRepo.GetMostUsedPromotion()
	assert.NoError(t, err, "Error fetching most used promotion from MySQL")

	switch p := mostUsed.(type) {
	case entities.DiscountEntitySQL:
		log.Fatalf("Error")
	case entities.FinancingEntitySQL:
		assert.Equal(t, p.PromotionEntitySQL.Code, "PV20241001")
	default:
		log.Fatalf("Error")
	}

	// ------ NoSQL (MongoDB) ------
	noSQLPromotionRepo := non_relational_repository.NewPromotionNonRelationalRepository(NoSQLDatabase)
	mostUsedMongo, err := noSQLPromotionRepo.GetMostUsedPromotion()

	assert.NoError(t, err, "Error fetching most used promotion from MongoDB")

	switch p := mostUsedMongo.(type) {
	case entities.DiscountEntityNonSQL:
		assert.Equal(t, p.PromotionEntity.Code, "WINTERSALE2024")
	case entities.FinancingEntityNonSQL:
		log.Fatal("This should not be a financing promotion")
	default:
		log.Fatalf("This should have been a discount promotion or a financing promotion")
	}
}

// ---------------------------------------------------
// --------------  STORE TESTS -----------------------
// ---------------------------------------------------

func TestStoreGetStoreWithHighestRevenueByMonth(t *testing.T) {
	storeRepo := relational_repository.NewStoreRelationalRepository(SQLDatabase)

	result, err := storeRepo.GetStoreWithHighestRevenueByMonth(10, 2024)
	assert.NoError(t, err, "Error fetching store with highest revenue by month from MySQL")

	assert.Equal(t, result.Cuit, "30-15066778-9")
	assert.Equal(t, result.Name, "Store O")

	// ------ NoSQL (MongoDB) ------
	noSQLStoreRepo := non_relational_repository.NewStoreNonRelationalRepository(NoSQLDatabase)
	resultMongo, err := noSQLStoreRepo.GetStoreWithHighestRevenueByMonth(10, 2024)

	assert.NoError(t, err, "Error fetching store with highest revenue by month from MongoDB")

	assert.Equal(t, resultMongo.Cuit, "30-12345678-9")
	assert.Equal(t, resultMongo.Name, "Store A")
}
