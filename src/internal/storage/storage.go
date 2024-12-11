package storage

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

// IStorage is the interface that wraps the basic storage methods.
// It includes methods to interact with the bank, card, customer, payment summary, purchase, and quota entities.
// The methods are not implemented in this interface and are delegated to the specific storage implementation, for
// relational databases, NoSQL databases, or other storage systems.
type IStorage interface {
	/* Bank */
	AddFinancingPromotionToBank(promotionFinancing models.Financing) error
	ExtendFinancingPromotionValidity(code string, newDate time.Time) error
	ExtendDiscountPromotionValidity(code string, newDate time.Time) error
	DeleteFinancingPromotion(code string) error
	DeleteDiscountPromotion(code string) error
	GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error)
	/* Card */

	/* Customer */

	/* PaymentSummary */

	/* Purchase */

	/* Quota */
}
