package models

// Quota represents an installment (or quota) for a purchase that is being paid in monthly payments.
type Quota struct {
	Number int     `json:"number"`
	Price  float64 `json:"price"`
	Month  string  `json:"month"`
	Year   string  `json:"year"`
}
