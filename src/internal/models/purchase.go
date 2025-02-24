/*
 * Payment Registration System - Purchase Model
 * -------------------------------------------
 * This file defines the data model for a purchase, representing a financial transaction made by a customer using a card.
 * It includes details like the store, initial and final amounts, and the type of purchase.
 *
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package models

// Purchase represents a financial transaction made at a store.
//
//	@Summary		Purchase model
//	@Description	Contains details about a purchase, including the store, amount, and type of purchase.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Purchase struct {
	PaymentVoucher string       `json:"payment_voucher" example:"VCHR-202502"` // Unique identifier for the purchase
	Store          string       `json:"store" example:"ElectroStore"`          // Name of the store where the purchase was made
	CuitStore      string       `json:"cuit_store" example:"30-98765432-1"`    // Unique tax identification code (CUIT) of the store
	Amount         float64      `json:"amount" example:"1500.75"`              // Initial purchase amount before any adjustments
	FinalAmount    float64      `json:"final_amount" example:"1400.00"`        // Final amount after discounts or interest
	PurchaseType   PurchaseType `json:"purchase_type" example:"0"`             // Type of purchase (single payment or installments)
}

// PurchaseSinglePayment represents a single-payment purchase.
//
//	@Summary		PurchaseSinglePayment model
//	@Description	Contains details about a single-payment purchase, including applicable store discounts.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type PurchaseSinglePayment struct {
	Purchase
	StoreDiscount float64 `json:"store_discount" example:"5.0"` // Discount applied to the purchase
}

// PurchaseMonthlyPayment represents a purchase paid in monthly installments.
//
//	@Summary		PurchaseMonthlyPayment model
//	@Description	Contains details about a monthly installment purchase, including interest rate, number of quotas, and quota breakdown.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type PurchaseMonthlyPayment struct {
	Purchase
	Interest       float64 `json:"interest" example:"3.5"`        // Interest rate applied to the purchase
	NumberOfQuotas int     `json:"number_of_quotas" example:"12"` // Number of monthly installments
	Quota          []Quota `json:"quota"`                         // Breakdown of installment payments
}

// PurchaseType represents the type of a purchase, either single payment or monthly payments.
//
//	@Summary		PurchaseType model
//	@Description	Enum representing whether a purchase is a single payment or paid in multiple installments.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type PurchaseType int

const (
	// SinglePayment represents a purchase that is paid in a single installment.
	SinglePayment PurchaseType = iota

	// MonthlyPayments represents a purchase that is paid in multiple installments.
	MonthlyPayments
)

// String returns the string representation of the PurchaseType.
func (p PurchaseType) String() string {
	return [...]string{"SinglePayment", "MonthlyPayments"}[p]
}
