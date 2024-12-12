package storage

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

// IBankStorage is the interface that defines methods related to bank operations,
// such as adding, deleting, and extending promotions, and fetching customer counts.
type IBankStorage interface {
	// AddFinancingPromotionToBank adds a financing promotion to the bank.
	AddFinancingPromotionToBank(promotionFinancing models.Financing) error
	// ExtendFinancingPromotionValidity extends the validity of a financing promotion.
	ExtendFinancingPromotionValidity(code string, newDate time.Time) error
	// ExtendDiscountPromotionValidity extends the validity of a discount promotion.
	ExtendDiscountPromotionValidity(code string, newDate time.Time) error
	// DeleteFinancingPromotion deletes a financing promotion by its code.
	DeleteFinancingPromotion(code string) error
	// DeleteDiscountPromotion deletes a discount promotion by its code.
	DeleteDiscountPromotion(code string) error
	// GetBankCustomerCounts retrieves the count of customers for each bank.
	GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error)
}

// ICardStorage is the interface that defines methods related to card operations,
// such as retrieving payment summaries, card expiration data, and purchases.
type ICardStorage interface {
	// GetPaymentSummary retrieves the payment summary for a card.
	GetPaymentSummary(cardNumber string, month int, year int) (*models.PaymentSummary, error)
	// GetCardsExpiringInNext30Days retrieves cards that will expire in the next 30 days.
	GetCardsExpiringInNext30Days(day int, month int, year int) (*[]models.Card, error)
	// GetPurchaseMonthly retrieves the monthly purchase details for a card.
	GetPurchaseMonthly(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseMonthlyPayment, error)
	// GetPurchaseSingle retrieves the single purchase details for a card.
	GetPurchaseSingle(cuit string, finalAmount float64, paymentVoucher string) (*models.PurchaseSinglePayment, error)
	// GetTop10CardsByPurchases retrieves the top 10 cards by purchases.
	GetTop10CardsByPurchases() (*[]models.Card, error)
}

// IPromotionStorage is the interface that defines methods related to promotion operations,
// such as retrieving available promotions and the most used promotion.
type IPromotionStorage interface {
	// GetAvailablePromotionsByStoreAndDateRange retrieves available promotions for a store within a date range.
	GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]models.Financing, *[]models.Discount, error)
	// GetMostUsedPromotion retrieves the most used promotion.
	GetMostUsedPromotion() (interface{}, error)
}

// IStoreStorage is the interface that defines methods related to store operations,
// such as retrieving the store with the highest revenue in a given month.
type IStoreStorage interface {
	// GetStoreWithHighestRevenueByMonth retrieves the store with the highest revenue in a specific month.
	GetStoreWithHighestRevenueByMonth(month int, year int) (models.StoreDTO, error)
}
