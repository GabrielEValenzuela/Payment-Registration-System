package models

import (
	"time"
)

// PaymentSummary represents a summary of payments for a card in a specific month and year.
// It includes payment due dates and increment rates for late payments.
type PaymentSummary struct {
	Code                string    `json:"code"`
	Month               int       `json:"month"`
	Year                int       `json:"year"`
	FirstExpiration     time.Time `json:"first_expiration"`
	SecondExpiration    time.Time `json:"second_expiration"`
	SurchargePercentage float64   `json:"surcharge_percentage"`
	TotalPrice          float64   `json:"total_price"`
	MonthlyPayments     []int     `json:"monthly_payments"`
	SinglePayments      []int     `json:"single_payments"`
	Card                []int     `json:"card"`
}
