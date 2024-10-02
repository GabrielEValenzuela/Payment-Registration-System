package customer

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
)

// Customer represents a person who is a client of one or more banks.
// A customer can have multiple cards and be associated with different banks.
type Customer struct {
	CompleteName string      `json:"complete_name"`
	Dni          string      `json:"dni"`
	Cuit         string      `json:"cuit"`
	Address      string      `json:"address"`
	Telephone    string      `json:"telephone"`
	EntryDate    time.Time   `json:"entry_date"`
	BanksIds     []int       `json:"backs_ids"`
	Cards        []card.Card `json:"cards"`
}
