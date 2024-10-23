package models

import (
	"time"
)

// Card represents a credit or debit card issued by a bank to a customer.
// A card is linked to a specific bank and customer, and it records purchases made with the card.
type Card struct {
	Number                  string                   `json:"number"`
	Ccv                     string                   `json:"ccv"`
	CardholderNameInCard    string                   `json:"cardholdername_in_card"`
	Since                   time.Time                `json:"since"`
	ExpirationDate          time.Time                `json:"expiration_date"`
	Bank                    Bank                     `json:"bank"`
	PurchaseMonthlyPayments []PurchaseMonthlyPayment `json:"purchase_monthly_payments"`
	PurchaseSinglePayments  []PurchaseSinglePayment  `json:"purchase_single_payment"`
}
