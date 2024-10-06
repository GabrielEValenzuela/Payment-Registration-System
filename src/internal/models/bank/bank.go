package bank

// Bank represents a financial institution that holds customers and issues cards.
type Bank struct {
	Name         string `json:"name"`
	Cuit         string `json:"cuit"`
	Address      string `json:"address"`
	Telephone    string `json:"telephone"`
	CustomersIds []int  `json:"customers_ids"`
}
