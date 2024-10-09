package gorm

import (
	"log"

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

	// Return the DB instance
	return db, nil
}
