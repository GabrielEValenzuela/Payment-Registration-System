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

// PaymentSummary represents a summary of payments for a specific period.
//
//	@Summary		Payment summary model
//	@Description	Contains details about a payment summary, including expiration dates, surcharge percentage, total price, and payment breakdown.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type PaymentSummary struct {
	Code                string                   `json:"code" example:"PAY-202502"`                        // Unique code identifying the payment summary
	Month               int                      `json:"month" example:"2"`                                // Month of the payment summary
	Year                int                      `json:"year" example:"2025"`                              // Year of the payment summary
	FirstExpiration     time.Time                `json:"first_expiration" example:"2025-02-10T00:00:00Z"`  // First expiration date
	SecondExpiration    time.Time                `json:"second_expiration" example:"2025-02-20T00:00:00Z"` // Second expiration date
	SurchargePercentage float64                  `json:"surcharge_percentage" example:"5.0"`               // Surcharge percentage applied after first expiration
	TotalPrice          float64                  `json:"total_price" example:"1500.75"`                    // Total price to be paid
	MonthlyPayments     []PurchaseMonthlyPayment `json:"monthly_payments"`                                 // List of monthly installment payments
	SinglePayments      []PurchaseSinglePayment  `json:"single_payments"`                                  // List of single-payment transactions
	Card                Card                     `json:"card"`                                             // Card used for the payment
}
