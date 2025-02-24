/*
 * Payment Registration System - Card Models
 * ----------------------------------------
 * This file defines the data models related to credit and debit cards,
 * including their associations with banks and recorded purchases.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package models

import (
	"time"
)

// Card represents a payment card issued by a bank.
//
//	@Summary		Card model
//	@Description	Contains details about a payment card, including its number, security code, cardholder name, issuance and expiration dates, associated bank, and purchase transactions.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Card struct {
	Number                  string                   `json:"number" example:"1234-5678-9012-3456"`           // Unique card number
	Ccv                     string                   `json:"ccv" example:"123"`                              // Card verification code
	CardholderNameInCard    string                   `json:"cardholdername_in_card" example:"John Doe"`      // Name as printed on the card
	Since                   time.Time                `json:"since" example:"2020-01-01T00:00:00Z"`           // Issuance date of the card
	ExpirationDate          time.Time                `json:"expiration_date" example:"2025-12-31T23:59:59Z"` // Expiration date of the card
	Bank                    Bank                     `json:"bank"`                                           // Issuing bank details
	PurchaseMonthlyPayments []PurchaseMonthlyPayment `json:"purchase_monthly_payments"`                      // Monthly installment payments
	PurchaseSinglePayments  []PurchaseSinglePayment  `json:"purchase_single_payment"`                        // Single-payment transactions
}
