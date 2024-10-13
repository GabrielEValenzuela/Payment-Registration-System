package purchase

import "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/quota"

// Purchase represents a financial transaction made by a customer using a card.
// It includes details like the store, initial and final amounts, and the type of purchase.
type Purchase struct {
	PaymentVoucher string       `json:"payment_voucher"`
	Store          string       `json:"store"`
	CuitStore      string       `json:"cuit_store"`
	Amount         float64      `json:"amount"`
	FinalAmount    float64      `json:"final_amount"`
	PurchaseType   PurchaseType `json:"purchase_type"` // Enum type
}

// PurchaseSinglePayment represents a one-time purchase made by a customer.
// It includes a discount percentage that might be applied at the store.
type PurchaseSinglePayment struct {
	Purchase
	StoreDiscount float64 `json:"store_discount"`
}

// PurchaseMonthlyPayments represents a purchase made in multiple installments.
// It includes an interest percentage and the number of installments.
type PurchaseMonthlyPayments struct {
	Purchase
	Interest       float64       `json:"interest"`
	NumberOfQuotas int           `json:"number_of_quotas"`
	Quota          []quota.Quota `json:"quota"`
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
