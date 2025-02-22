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

/*
 * Purchase
 * ----------------------------------------
 * Represents a purchase made by a customer using a card.
 *
 * Fields:
 * - PaymentVoucher (string): The payment voucher or receipt number.
 * - Store (string): The name of the store where the purchase was made.
 * - CuitStore (string): The unique tax identification code for the store.
 * - Amount (float64): The initial amount of the purchase.
 * - FinalAmount (float64): The final amount after any discounts or interest.
 * - PurchaseType (PurchaseType): The type of purchase, either single payment or monthly payments.
 */
type Purchase struct {
	PaymentVoucher string       `json:"payment_voucher"`
	Store          string       `json:"store"`
	CuitStore      string       `json:"cuit_store"`
	Amount         float64      `json:"amount"`
	FinalAmount    float64      `json:"final_amount"`
	PurchaseType   PurchaseType `json:"purchase_type"` // Enum type
}

/*
 * PurchaseSinglePayment
 * ----------------------------------------
 * Represents a purchase made in a single payment.
 *
 * Fields:
 * - Purchase (Purchase): The base purchase details.
 * - StoreDiscount (float64): The discount applied by the store.
 */
type PurchaseSinglePayment struct {
	Purchase
	StoreDiscount float64 `json:"store_discount"`
}

/*
 * PurchaseMonthlyPayment
 * ----------------------------------------
 * Represents a purchase made with monthly payments.
 *
 * Fields:
 * - Purchase (Purchase): The base purchase details.
 * - Interest (float64): The interest rate applied to the purchase.
 * - NumberOfQuotas (int): The number of installments for the purchase.
 * - Quota ([]Quota): List of installment details.
 */
type PurchaseMonthlyPayment struct {
	Purchase
	Interest       float64 `json:"interest"`
	NumberOfQuotas int     `json:"number_of_quotas"`
	Quota          []Quota `json:"quota"`
}

// PurchaseType represents the type of a purchase, either single payment or monthly payments.
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
