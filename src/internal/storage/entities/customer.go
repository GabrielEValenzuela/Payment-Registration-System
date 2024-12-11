package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerEntityNonSQL struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`        // MongoDB primary key
	CompleteName string               `bson:"complete_name"`        // Full name of the customer
	Dni          string               `bson:"dni"`                  // Unique identifier
	Cuit         string               `bson:"cuit"`                 // Unique identifier
	Address      string               `bson:"address,omitempty"`    // Optional address
	Telephone    string               `bson:"telephone,omitempty"`  // Optional telephone
	EntryDate    time.Time            `bson:"entry_date"`           // Date of customer entry
	Banks        []primitive.ObjectID `bson:"banks,omitempty"`      // References to related banks (IDs)
	Cards        []primitive.ObjectID `bson:"cards,omitempty"`      // References to related cards (IDs)
	CreatedAt    time.Time            `bson:"created_at,omitempty"` // Creation timestamp
	UpdatedAt    time.Time            `bson:"updated_at,omitempty"` // Update timestamp
}

// Bank represents a financial institution that holds customers and issues cards.
type CustomerEntitySQL struct {
	ID           uint            `gorm:"primaryKey;autoIncrement"`
	CompleteName string          `gorm:"size:255;not null"`
	Dni          string          `gorm:"size:20;unique;not null"`
	Cuit         string          `gorm:"size:20;unique;not null"`
	Address      string          `gorm:"size:255"`
	Telephone    string          `gorm:"size:50"`
	EntryDate    time.Time       `gorm:"not null"`
	Banks        []BankEntitySQL `gorm:"many2many:CUSTOMERS_BANKS"`
	Cards        []CardEntitySQL `gorm:"foreignKey:CustomerID"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime"`
}

func (CustomerEntitySQL) TableName() string {
	return "CUSTOMERS"
}

// ------------ Mappers ------------	//

// Take a model and convert it to a BankEntity for relational storage
func ToCustomerEntityRelational(customer *models.Customer) *CustomerEntitySQL {
	return &CustomerEntitySQL{
		CompleteName: customer.CompleteName,
		Dni:          customer.Dni,
		Cuit:         customer.Cuit,
		Address:      customer.Address,
		Telephone:    customer.Telephone,
		EntryDate:    customer.EntryDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Take a model and convert it to a BankEntity for non-relational storage
func ToCustomerEntityNonRelational(customer *models.Customer) *CustomerEntityNonSQL {
	return &CustomerEntityNonSQL{
		CompleteName: customer.CompleteName,
		Dni:          customer.Dni,
		Cuit:         customer.Cuit,
		Address:      customer.Address,
		Telephone:    customer.Telephone,
		EntryDate:    customer.EntryDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CustomerModel a Customer mapper (si necesitas convertir de nuevo)
func ToCustomer[T any](customerEntity *T) *models.Customer {
	switch v := any(customerEntity).(type) {
	case *CustomerEntitySQL:
		return &models.Customer{
			CompleteName: v.CompleteName,
			Dni:          v.Dni,
			Cuit:         v.Cuit,
			Address:      v.Address,
			Telephone:    v.Telephone,
			EntryDate:    v.EntryDate,
		}
	case *CustomerEntityNonSQL:
		return &models.Customer{
			CompleteName: v.CompleteName,
			Dni:          v.Dni,
			Cuit:         v.Cuit,
			Address:      v.Address,
			Telephone:    v.Telephone,
			EntryDate:    v.EntryDate,
		}
	}
	return nil
}
