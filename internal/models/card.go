package models

// Card represents a credit or debit card issued by a bank to a customer.
// A card is linked to a specific bank and customer, and it records purchases made with the card.
type Card struct {
	ID         int        `json:"id"`
	CardNumber string     `json:"card_number"`
	Bank       Bank       `json:"bank"`
	Customer   Customer   `json:"customer"`
	Purchases  []Purchase `json:"purchases"`
}
