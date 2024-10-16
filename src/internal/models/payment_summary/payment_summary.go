package payment_summary

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/purchase"
)

// PaymentSummary represents a summary of payments for a card in a specific month and year.
// It includes payment due dates and increment rates for late payments.
type PaymentSummary struct {
	Code                string                            `json:"code"`
	Month               int                               `json:"month"`
	Year                int                               `json:"year"`
	FirstExpiration     time.Time                         `json:"first_expiration"`
	SecondExpiration    time.Time                         `json:"second_expiration"`
	SurchargePercentage float64                           `json:"surcharge_percentage"`
	TotalPrice          float64                           `json:"total_price"`
	MonthlyPayments     []purchase.PurchaseMonthlyPayment `json:"monthly_payments"`
	SinglePayments      []purchase.PurchaseSinglePayment  `json:"single_payments"`
	Card                card.Card                         `json:"card"`
}
