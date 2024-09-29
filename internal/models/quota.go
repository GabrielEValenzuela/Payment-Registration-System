package models

// Quota represents an installment (or quota) for a purchase that is being paid in monthly payments.
type Quota struct {
	ID       int                     `json:"id"`
	Amount   float64                 `json:"amount"`
	Purchase PurchaseMonthlyPayments `json:"purchase"`
}
