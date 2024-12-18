package services

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
)

// StoreService defines the interface for store-related operations.
// This service abstracts business logic and data layer interactions,
// providing a clear contract for retrieving store-related data like revenue and store-specific information.
type StoreService interface {
	// GetStoreWithHighestRevenueByMonth retrieves the store with the highest revenue in a specific month and year.
	// Parameters:
	// - month: The month for which to retrieve the highest revenue store.
	// - year: The year for which to retrieve the highest revenue store.
	// Returns:
	// - models.StoreDTO: The store with the highest revenue.
	// - error: An error if the operation fails, otherwise nil.
	GetStoreWithHighestRevenueByMonth(month int, year int) (models.StoreDTO, error)
}

// storeService is a concrete implementation of the StoreService interface.
// It uses a repository (IStoreStorage) to perform data operations.
type storeService struct {
	repo storage.IStoreStorage
}

// NewStoreService creates and initializes a new StoreService instance.
// Parameters:
// - repo: An IStoreStorage repository interface for interacting with the data layer.
// Returns:
// - StoreService: A new instance of the service struct implementing the StoreService interface.
func NewStoreService(repo storage.IStoreStorage) StoreService {
	return &storeService{
		repo: repo,
	}
}

// GetStoreWithHighestRevenueByMonth retrieves the store with the highest revenue in a specific month and year.
func (s *storeService) GetStoreWithHighestRevenueByMonth(month int, year int) (models.StoreDTO, error) {
	return s.repo.GetStoreWithHighestRevenueByMonth(month, year)
}
