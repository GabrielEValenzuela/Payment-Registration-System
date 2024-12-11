package relational

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

// NewMySQLDB creates a new connection to the MySQL database and initializes the schema.
func NewMySQLDB(dsn string, cleanDB bool) (*gorm.DB, error) {
	newLogger := gorm_logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gorm_logger.Config{
			SlowThreshold:             time.Second,        // Slow SQL threshold
			LogLevel:                  gorm_logger.Silent, // Log level
			IgnoreRecordNotFoundError: false,               // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,               // Don't include params in the SQL log
			Colorful:                  true,              // Disable color
		},
	)
	// Open a new GORM connection using the MySQL driver
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	// Initialize the database schema
	if err := initSQLDB(db, cleanDB); err != nil {
		return nil, fmt.Errorf("failed to initialize MySQL database: %w", err)
	}

	return db, nil
}

// CloseDB gracefully closes the database connection.
func CloseDB(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// initSQLDB initializes the database schema.
// If `cleanDB` is true, it will drop existing tables before migrating.
func initSQLDB(database *gorm.DB, cleanDB bool) error {
	if cleanDB {
		logger.Info("Cleaning the database: dropping existing tables...")
		// Retrieve the list of tables
		var tables []string
		if err := database.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
			return fmt.Errorf("failed to retrieve table list: %w", err)
		}

		// Drop each table
		for _, table := range tables {
			if err := database.Migrator().DropTable(table); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}
	}

	logger.Info("Migrating database schema...")
	// Automate the creation of tables based on entity definitions
	if err := database.AutoMigrate(
		&entities.BankEntitySQL{},
		&entities.CustomerEntitySQL{},
		&entities.CardEntitySQL{},
		&entities.QuotaEntitySQL{},
		&entities.PurchaseMonthlyPaymentsEntitySQL{},
		&entities.PurchaseSinglePaymentEntitySQL{},
		&entities.DiscountEntitySQL{},
		&entities.FinancingEntitySQL{},
		&entities.PaymentSummaryEntitySQL{},
	); err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}

	logger.Info("Database schema initialized successfully.")
	return nil
}
