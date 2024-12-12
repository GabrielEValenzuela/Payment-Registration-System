package services

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
)

// CardService defines the interface for card-related operations.
// This service abstracts business logic and data layer interactions,
// providing a clear contract for managing card operations like payment summaries, purchases, and expiring cards.
type CardService interface {
	// GetPaymentSummary retrieves the payment summary for a card.
	// Parameters:
	// - cardNumber: The card number for which the payment summary is required.
	// - month: The month for the payment summary.
	// - year: The year for the payment summary.
	// Returns:
	// - *models.PaymentSummary: A PaymentSummary object containing the payment details for the specified card.
	// - error: An error if the operation fails, otherwise nil.
	GetPaymentSummary(cardNumber string, month int, year int) (*models.PaymentSummary, error)

	// GetCardsExpiringInNext30Days retrieves the cards that will expire in the next 30 days.
	// Parameters:
	// - day: The current day.
	// - month: The current month.
	// - year: The current year.
	// Returns:
	// - *[]models.Card: A slice of Card objects representing the cards expiring in the next 30 days.
	// - error: An error if the operation fails, otherwise nil.
	GetCardsExpiringInNext30Days(day int, month int, year int) (*[]models.Card, error)

	// GetPurchaseMonthly retrieves the monthly purchase details for a card.
	// Parameters:
	// - cuit: The CUIT of the card holder.
	// - finalAmount: The final amount for the purchase.
	// - paymentVoucher: The payment voucher ID.
	// Returns:
	// - *models.PurchaseMonthlyPayment: A PurchaseMonthlyPayment object containing the purchase details for the card.
	// - error: An error if the operation fails, otherwise nil.
	GetPurchaseMonthly(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseMonthlyPayment, error)

	// GetPurchaseSingle retrieves the single purchase details for a card.
	// Parameters:
	// - cuit: The CUIT of the card holder.
	// - finalAmount: The final amount for the purchase.
	// - paymentVoucher: The payment voucher ID.
	// Returns:
	// - *models.PurchaseSinglePayment: A PurchaseSinglePayment object containing the purchase details for the card.
	// - error: An error if the operation fails, otherwise nil.
	GetPurchaseSingle(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseSinglePayment, error)

	// GetTop10CardsByPurchases retrieves the top 10 cards by the number of purchases.
	// Returns:
	// - *[]models.Card: A slice of Card objects representing the top 10 cards by purchases.
	// - error: An error if the operation fails, otherwise nil.
	GetTop10CardsByPurchases() (*[]models.Card, error)
}

// service is a concrete implementation of the CardService interface.
// It uses a repository (ICardStorage) to perform data operations.
type cardService struct {
	repo storage.ICardStorage
}

// NewCardService creates and initializes a new CardService instance.
// Parameters:
// - repo: An ICardStorage repository interface for interacting with the data layer.
// Returns:
// - CardService: A new instance of the service struct implementing the CardService interface.
func NewCardService(repo storage.ICardStorage) CardService {
	return &cardService{
		repo: repo,
	}
}

// GetPaymentSummary retrieves the payment summary for a card.
func (s *cardService) GetPaymentSummary(cardNumber string, month int, year int) (*models.PaymentSummary, error) {
	return s.repo.GetPaymentSummary(cardNumber, month, year)
}

// GetCardsExpiringInNext30Days retrieves the cards that will expire in the next 30 days.
func (s *cardService) GetCardsExpiringInNext30Days(day int, month int, year int) (*[]models.Card, error) {
	return s.repo.GetCardsExpiringInNext30Days(day, month, year)
}

// GetPurchaseMonthly retrieves the monthly purchase details for a card.
func (s *cardService) GetPurchaseMonthly(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseMonthlyPayment, error) {
	return s.repo.GetPurchaseMonthly(cuit, finalAmount, paymentVoucher)
}

// GetPurchaseSingle retrieves the single purchase details for a card.
func (s *cardService) GetPurchaseSingle(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseSinglePayment, error) {
	return s.repo.GetPurchaseSingle(cuit, finalAmount, paymentVoucher)
}

// GetTop10CardsByPurchases retrieves the top 10 cards by purchases.
func (s *cardService) GetTop10CardsByPurchases() (*[]models.Card, error) {
	return s.repo.GetTop10CardsByPurchases()
}
