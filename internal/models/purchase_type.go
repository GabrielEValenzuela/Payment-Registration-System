package models

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
