package store

import "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"

type Repository interface {
	GetStoreWithHighestRevenueByMonth(month int, year int) (models.StoreDTO, error)
}
