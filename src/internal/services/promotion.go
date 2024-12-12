package services

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
)

// PromotionService defines the interface for promotion-related operations.
// This service abstracts business logic and data layer interactions,
// providing a clear contract for managing promotions like financing, discounts, and store-specific promotions.
type PromotionService interface {
	// GetAvailablePromotionsByStoreAndDateRange retrieves available promotions by store and date range.
	// Parameters:
	// - cuit: The CUIT of the store.
	// - startDate: The start date of the promotion period.
	// - endDate: The end date of the promotion period.
	// Returns:
	// - *[]models.Financing: A slice of available financing promotions.
	// - *[]models.Discount: A slice of available discount promotions.
	// - error: An error if the operation fails, otherwise nil.
	GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]models.Financing, *[]models.Discount, error)

	// GetMostUsedPromotion retrieves the most used promotion.
	// Returns:
	// - interface{}: The most used promotion.
	// - error: An error if the operation fails, otherwise nil.
	GetMostUsedPromotion() (interface{}, error)
}

// promotionService is a concrete implementation of the PromotionService interface.
// It uses a repository (IPromotionStorage) to perform data operations.
type promotionService struct {
	repo storage.IPromotionStorage
}

// NewPromotionService creates and initializes a new PromotionService instance.
// Parameters:
// - repo: An IPromotionStorage repository interface for interacting with the data layer.
// Returns:
// - PromotionService: A new instance of the service struct implementing the PromotionService interface.
func NewPromotionService(repo storage.IPromotionStorage) PromotionService {
	return &promotionService{
		repo: repo,
	}
}

// GetAvailablePromotionsByStoreAndDateRange retrieves available promotions by store and date range.
func (s *promotionService) GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]models.Financing, *[]models.Discount, error) {
	return s.repo.GetAvailablePromotionsByStoreAndDateRange(cuit, startDate, endDate)
}

// GetMostUsedPromotion retrieves the most used promotion.
func (s *promotionService) GetMostUsedPromotion() (interface{}, error) {
	return s.repo.GetMostUsedPromotion()
}
