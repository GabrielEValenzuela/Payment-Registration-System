package models

// Bank represents a financial institution that holds customers and issues cards.
type Bank struct {
	Name      string     `json:"name"`
	Cuit      string     `json:"cuit"`
	Address   string     `json:"address"`
	Telephone string     `json:"telephone"`
	Members   []Customer `json:"customers_ids"`
}

type BankCustomerCountDTO struct {
	BankCuit      string `json:"bank_cuit"`
	BankName      string `json:"bank_name"`
	CustomerCount int    `json:"customer_count"`
}
