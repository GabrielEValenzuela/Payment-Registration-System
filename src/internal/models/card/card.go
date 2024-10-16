package card

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/bank"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/purchase"
)

// Card represents a credit or debit card issued by a bank to a customer.
// A card is linked to a specific bank and customer, and it records purchases made with the card.
type Card struct {
	Number                  string                            `json:"number"`
	Ccv                     string                            `json:"ccv"`
	CardholderNameInCard    string                            `json:"cardholdername_in_card"`
	Since                   time.Time                         `json:"since"`
	ExpirationDate          time.Time                         `json:"expiration_date"`
	Bank                    bank.Bank                         `json:"bank"`
	PurchaseMonthlyPayments []purchase.PurchaseMonthlyPayment `json:"purchase_monthly_payments"`
	PurchaseSinglePayments  []purchase.PurchaseSinglePayment  `json:"purchase_single_payment"`
}
