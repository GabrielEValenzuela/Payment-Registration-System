package promotion

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/bank"
)

// Promotion represents a special offer provided by a bank to customers.
// It applies to specific stores and is valid for a certain period of time.
type Promotion struct {
	Code              string    `json:"code"`
	PromotionTitle    string    `json:"promotion_title"`
	NameStore         string    `json:"name_store"`
	CuitStore         string    `json:"cuit_store"`
	ValidityStartDate time.Time `json:"validity_start_date"`
	ValidityEndDate   time.Time `json:"validity_end_date"`
	Comments          string    `json:"comments"`
	Bank              bank.Bank `json:"bank"`
}

// Discount represents a type of promotion that applies a percentage discount to purchases.
// Optionally, it may have a maximum discount amount.
type Discount struct {
	Promotion
	DiscountPercentage float64 `json:"DiscountPercentage"`
	PriceCap           float64 `json:"PriceCap"`
	OnlyCash           bool    `json:"OnlyCash"`
}

// Financing represents a promotion that offers installment payment options with specific interest rates.
type Financing struct {
	Promotion
	NumberOfQuotas int     `json:"number_of_quotas"`
	Interest       float64 `json:"interest"`
}
