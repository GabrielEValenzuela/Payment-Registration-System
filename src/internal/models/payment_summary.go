/*
* Payment Registration System - Payment Summary Model
* ----------------------------------------------------
* This file defines the data model for a payment summary, representing a summary of payments
* for a card in a specific month and year. It includes payment due dates and increment rates
* for late payments.
*
* Created: Oct. 19, 2024
* License: GNU General Public License v3.0
 */
package models

import (
	"time"
)

/* PurchaseMonthlyPayment
 * ----------------------------------------
 * Represents a monthly installment payment made with a card.
 *
 * Fields:
 * - Code (string): The unique identifier of the payment.
 * - Month (int): The month of the payment.
 * - Year (int): The year of the payment.
 * - FirstExpiration (time.Time): The first due date for the payment.
 * - SecondExpiration (time.Time): The second due date for the payment.
 * - SurchargePercentage (float64): The percentage of surcharge for late payments.
 * - TotalPrice (float64): The total price of the payment.
 * - MonthlyPayments([]PurchaseMonthlyPayment): List of monthly installment payments linked to this card.
 * - SinglePayments([]PurchaseSinglePayment): List of single-payment purchases made with this card.
 * - Card (Card): The card associated with this payment.
 */
type PaymentSummary struct {
	Code                string                   `json:"code"`
	Month               int                      `json:"month"`
	Year                int                      `json:"year"`
	FirstExpiration     time.Time                `json:"first_expiration"`
	SecondExpiration    time.Time                `json:"second_expiration"`
	SurchargePercentage float64                  `json:"surcharge_percentage"`
	TotalPrice          float64                  `json:"total_price"`
	MonthlyPayments     []PurchaseMonthlyPayment `json:"monthly_payments"`
	SinglePayments      []PurchaseSinglePayment  `json:"single_payments"`
	Card                Card                     `json:"card"`
}
