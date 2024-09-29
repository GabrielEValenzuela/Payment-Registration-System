package models

// Purchase represents a financial transaction made by a customer using a card.
// It includes details like the store, initial and final amounts, and the type of purchase.
type Purchase struct {
	ID            int          `json:"id"`
	StoreName     string       `json:"store_name"`
	InitialAmount float64      `json:"initial_amount"`
	FinalAmount   float64      `json:"final_amount"`
	PurchaseType  PurchaseType `json:"purchase_type"` // Enum type
	Card          Card         `json:"card"`
}

// PurchaseSinglePayment represents a one-time purchase made by a customer.
// It includes a discount percentage that might be applied at the store.
type PurchaseSinglePayment struct {
	Purchase
	DiscountPercentage float64 `json:"discount_percentage"`
}

// PurchaseMonthlyPayments represents a purchase made in multiple installments.
// It includes an interest percentage and the number of installments.
type PurchaseMonthlyPayments struct {
	Purchase
	InterestPercentage float64 `json:"interest_percentage"`
	Installments       int     `json:"installments"`
}
