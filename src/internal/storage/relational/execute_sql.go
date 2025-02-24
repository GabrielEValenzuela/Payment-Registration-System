package relational

import (
	"fmt"
	"os"
	"strings"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"gorm.io/gorm"
)

func ExecuteSQLFile(db *gorm.DB, filePath string) error {
	// Leer el archivo SQL utilizando os.ReadFile
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %v", err)
	}

	// Dividir el contenido del archivo en múltiples consultas si es necesario
	queries := strings.Split(string(sqlContent), ";\n")

	// Ejecutar cada consulta SQL
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

// ShouldInitializeData checks if the initialization script needs to be run again
func ShouldInitializeData(db *gorm.DB) (bool, error) {
	var bankEntity entities.BankEntitySQL
	cuit := "30-12345678-9"

	// Check if a bank with ID 1 and name "Santander" exists
	err := db.First(&bankEntity, "cuit = ?", cuit).Error

	if err != nil {
		// If no record is found, return true to indicate the script should run
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		// Return false and the error for other issues
		return false, err
	}

	// If the bank exists, no need to run the initialization script
	return false, nil
}
