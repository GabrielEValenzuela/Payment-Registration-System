package entities

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Store represents a store and its associated information.
type StoreNonSQL struct {
	ID          bson.ObjectID `bson:"_id,omitempty"` // MongoDB primary key
	StoreName   string        `bson:"store_name"`    // Name of the store
	CuitStore   string        `bson:"cuit_store"`    // Store CUIT
	TotalAmount float64       `bson:"total_amount"`  // Total amount associated with the store
}

// StoreSQL represents a store and its associated information.
type StoreSQL struct {
	Store       string
	CuitStore   string
	TotalAmount float64
}

func (StoreSQL) TableName() string {
	return "STORES"
}
