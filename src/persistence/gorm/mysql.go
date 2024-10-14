package gorm

import (
	"log"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQLDB creates a new connection to the MySQL database
func NewMySQLDB() (*gorm.DB, error) {
	// Define the DSN (Data Source Name) for MySQL connection
	dsn := "testuser:testpassword@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

	// Open a new GORM connection using the MySQL driver
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
		return nil, err
	}

	err = initSQLDB(db)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
		return nil, err
	}

	// Return the DB instance
	return db, nil
}

func CloseDB(database *gorm.DB) error {
	db, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
		return err
	}
	db.Close()
	return nil
}

func initSQLDB(database *gorm.DB) error {
	// Clean Database
	var tables []string
	database.Raw("SHOW TABLES").Scan(&tables)

	for _, table := range tables {
		database.Migrator().DropTable(table)
	}

	// Create tablets
	err := database.AutoMigrate(
		&entities.BankEntity{},
		&entities.CustomerEntity{},
		&entities.CardEntity{},
		&entities.QuotaEntity{},
		&entities.PurchaseMonthlyPaymentsEntity{},
		&entities.PurchaseSinglePaymentEntity{},
		&entities.DiscountEntity{},
		&entities.FinancingEntity{},
		&entities.PaymentSummaryEntity{},
	)

	return err
}
