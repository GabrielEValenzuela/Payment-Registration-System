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

/*
 * Card
 * ----------------------------------------
 * Represents a credit or debit card issued by a bank to a customer.
 * A card is linked to a specific bank and customer, and it records purchases made with the card.
 *
 * Fields:
 * - Number (string): Unique identifier of the card (e.g., 16-digit card number).
 * - Ccv (string): Card verification code (e.g., 3 or 4 digits).
 * - CardholderNameInCard (string): The name printed on the card.
 * - Since (time.Time): The date when the card was issued.
 * - ExpirationDate (time.Time): The date when the card expires.
 * - Bank (Bank): The issuing bank associated with this card.
 * - PurchaseMonthlyPayments ([]PurchaseMonthlyPayment): List of monthly installment purchases linked to this card.
 * - PurchaseSinglePayments ([]PurchaseSinglePayment): List of single-payment purchases made with this card.
 */
type Card struct {
	Number                  string                   `json:"number"`                    // Unique card number
	Ccv                     string                   `json:"ccv"`                       // Card verification code
	CardholderNameInCard    string                   `json:"cardholdername_in_card"`    // Name as printed on the card
	Since                   time.Time                `json:"since"`                     // Issuance date of the card
	ExpirationDate          time.Time                `json:"expiration_date"`           // Expiration date of the card
	Bank                    Bank                     `json:"bank"`                      // Issuing bank details
	PurchaseMonthlyPayments []PurchaseMonthlyPayment `json:"purchase_monthly_payments"` // Monthly installment payments
	PurchaseSinglePayments  []PurchaseSinglePayment  `json:"purchase_single_payment"`   // Single-payment transactions
}
