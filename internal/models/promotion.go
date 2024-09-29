package models

import "time"

// Promotion represents a special offer provided by a bank to customers.
// It applies to specific stores and is valid for a certain period of time.
type Promotion struct {
	ID        int       `json:"id"`
	CUIT      string    `json:"cuit"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Bank      Bank      `json:"bank"`
}

// Discount represents a type of promotion that applies a percentage discount to purchases.
// Optionally, it may have a maximum discount amount.
type Discount struct {
	Promotion
	DiscountPercentage float64 `json:"discount_percentage"`
	MaxDiscountAmount  float64 `json:"max_discount_amount,omitempty"`
}

// Financing represents a promotion that offers installment payment options with specific interest rates.
type Financing struct {
	Promotion
	Installments       int     `json:"installments"`
	InterestPercentage float64 `json:"interest_percentage"`
}
