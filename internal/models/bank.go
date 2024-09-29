package models

// Bank represents a financial institution that holds customers and issues cards.
type Bank struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Customers []Customer `json:"customers"`
}
