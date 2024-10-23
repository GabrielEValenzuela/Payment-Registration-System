package models

import (
	"time"
)

// Customer represents a person who is a client of one or more banks.
// A customer can have multiple cards and be associated with different banks.
type Customer struct {
	CompleteName string    `json:"complete_name" validate:"required"`
	Dni          string    `json:"dni" validate:"required"`
	Cuit         string    `json:"cuit" validate:"required"`
	Address      string    `json:"address" validate:"required"`
	Telephone    string    `json:"telephone" validate:"required"`
	EntryDate    time.Time `json:"entry_date" validate:"required"`
	BanksIds     []int     `json:"banks_ids" validate:"dive,required"`
	Cards        []int     `json:"cards" validate:"dive,required"`
}
