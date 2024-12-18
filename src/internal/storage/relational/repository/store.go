package relational_repository

import (
	"fmt"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
	"gorm.io/gorm"
)

type StoreRepositoryGORM struct {
	db *gorm.DB
}

// NewStoreRepository crea una nueva instancia de StoreRepository
func NewStoreRelationalRepository(db *gorm.DB) storage.IStoreStorage {
	return &StoreRepositoryGORM{db: db}
}

func (r *StoreRepositoryGORM) GetStoreWithHighestRevenueByMonth(month int, year int) (models.StoreDTO, error) {
	var result models.StoreDTO

	query := "SELECT store as name, cuit_store as cuit, SUM(final_amount) AS total_amount " +
		"FROM (" +
		"SELECT month.store, month.cuit_store, month.final_amount " +
		"FROM PURCHASE_MONTHLY_PAYMENTS month " +
		"WHERE MONTH(month.created_at) = %d AND YEAR(month.created_at) = %d " +
		"UNION ALL " +
		"SELECT single.store, single.cuit_store, single.final_amount " +
		"FROM PURCHASE_SINGLE_PAYMENTS single " +
		"WHERE MONTH(single.created_at) = %d AND YEAR(single.created_at) = %d" +
		") AS combined_payments " +
		"GROUP BY store, cuit_store " +
		"ORDER BY total_amount DESC;"

	formattedQuery := fmt.Sprintf(query, month, year, month, year)

	r.db.Raw(formattedQuery).Scan(&result)

	return result, nil
}
