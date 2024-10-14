package testresource

import (
	"fmt"
	"os"
	"strings"

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
