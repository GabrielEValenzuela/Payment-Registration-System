package services

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
)

// BankService defines the interface for bank-related operations.
// This service abstracts business logic and data layer interactions,
// providing a clear contract for managing banks and their associated promotions and customers.
type BankService interface {

	// AddFinancingPromotionToBank adds a new financing promotion to a specific bank.
	// Parameters:
	// - promotionFinancing: A Financing object containing the promotion details.
	// Returns:
	// - error: An error if the operation fails, otherwise nil.
	AddFinancingPromotionToBank(promotionFinancing models.Financing) error

	// ExtendFinancingPromotionValidity extends the validity period of a financing promotion.
	// Parameters:
	// - code: The unique identifier of the promotion.
	// - newDate: The new expiration date for the promotion.
	// Returns:
	// - error: An error if the operation fails, otherwise nil.
	ExtendFinancingPromotionValidity(code string, newDate time.Time) error

	// ExtendDiscountPromotionValidity extends the validity period of a discount promotion.
	// Parameters:
	// - code: The unique identifier of the promotion.
	// - newDate: The new expiration date for the promotion.
	// Returns:
	// - error: An error if the operation fails, otherwise nil.
	ExtendDiscountPromotionValidity(code string, newDate time.Time) error

	// DeleteFinancingPromotion logically deletes a financing promotion by marking it as inactive.
	// Parameters:
	// - code: The unique identifier of the promotion.
	// Returns:
	// - error: An error if the operation fails, otherwise nil.
	DeleteFinancingPromotion(code string) error

	// DeleteDiscountPromotion logically deletes a discount promotion by marking it as inactive.
	// Parameters:
	// - code: The unique identifier of the promotion.
	// Returns:
	// - error: An error if the operation fails, otherwise nil.
	DeleteDiscountPromotion(code string) error

	// GetBankCustomerCounts retrieves the count of customers associated with each bank.
	// Returns:
	// - []models.BankCustomerCountDTO: A slice of BankCustomerCountDTO containing the bank name, CUIT, and customer count.
	// - error: An error if the operation fails, otherwise nil.
	GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error)
}

// service is a concrete implementation of the BankService interface.
// It uses a repository (IStorage) to perform data operations.
type bankService struct {
	repo storage.IBankStorage
}

// NewBankService creates and initializes a new BankService instance.
// Parameters:
// - repo: An IStorage repository interface for interacting with the data layer.
// Returns:
// - BankService: A new instance of the service struct implementing the BankService interface.
func NewBankService(repo storage.IBankStorage) BankService {
	return &bankService{
		repo: repo,
	}
}

// AddFinancingPromotionToBank adds a new financing promotion to a specific bank.
func (s *bankService) AddFinancingPromotionToBank(promotionFinancing models.Financing) error {
	return s.repo.AddFinancingPromotionToBank(promotionFinancing)
}

// ExtendFinancingPromotionValidity extends the validity period of a financing promotion.
func (s *bankService) ExtendFinancingPromotionValidity(code string, newDate time.Time) error {
	return s.repo.ExtendFinancingPromotionValidity(code, newDate)
}

// ExtendDiscountPromotionValidity extends the validity period of a discount promotion.
func (s *bankService) ExtendDiscountPromotionValidity(code string, newDate time.Time) error {
	return s.repo.ExtendDiscountPromotionValidity(code, newDate)
}

// DeleteFinancingPromotion logically deletes a financing promotion by marking it as inactive.
func (s *bankService) DeleteFinancingPromotion(code string) error {
	return s.repo.DeleteFinancingPromotion(code)
}

// DeleteDiscountPromotion logically deletes a discount promotion by marking it as inactive.
func (s *bankService) DeleteDiscountPromotion(code string) error {
	return s.repo.DeleteDiscountPromotion(code)
}

// GetBankCustomerCounts retrieves the count of customers associated with each bank.
func (s *bankService) GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error) {
	return s.repo.GetBankCustomerCounts()
}
