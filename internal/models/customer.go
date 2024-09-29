package models

// Customer represents a person who is a client of one or more banks.
// A customer can have multiple cards and be associated with different banks.
type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Banks     []Bank `json:"banks"`
	Cards     []Card `json:"cards"`
}
