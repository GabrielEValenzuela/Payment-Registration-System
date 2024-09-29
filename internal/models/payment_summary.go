package models

import "time"

// PaymentSummary represents a summary of payments for a card in a specific month and year.
// It includes payment due dates and increment rates for late payments.
type PaymentSummary struct {
	ID             int                     `json:"id"`
	Month          int                     `json:"month"`
	Year           int                     `json:"year"`
	Card           Card                    `json:"card"`
	InternalCode   string                  `json:"internal_code"`
	DueDate1       time.Time               `json:"due_date_1"`
	DueDate2       time.Time               `json:"due_date_2"`
	IncrementRate1 float64                 `json:"increment_rate_1"`
	IncrementRate2 float64                 `json:"increment_rate_2"`
	Quotas         []Quota                 `json:"quotas"`
	SinglePayments []PurchaseSinglePayment `json:"single_payments"`
}
