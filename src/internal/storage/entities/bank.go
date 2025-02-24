/*
 * Payment Registration System - Bank Entity (SQL and NoSQL)
 * ---------------------------------------------------------
 *
 * Description: Bank entity represents a financial institution that holds customers and issues cards.
 * The entity is implemented in two ways: SQL and NoSQL.
 * The SQL implementation uses GORM and the NoSQL implementation uses MongoDB.
 * The entity has a one-to-many relationship with the Customer entity.
 *
 * Created: Dec. 11, 2024
 * License: GNU General Public License v3.0
 */

package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BankEntityNonSQL struct {
	ID        bson.ObjectID   `bson:"_id,omitempty"` // 🔥 Ensure `_id` exists
	Name      string          `bson:"name"`
	Cuit      string          `bson:"cuit"`
	Address   string          `bson:"address"`
	Telephone string          `bson:"telephone"`
	Customers []bson.ObjectID `bson:"customers,omitempty"`
	CreatedAt time.Time       `bson:"created_at,omitempty"`
	UpdatedAt time.Time       `bson:"updated_at,omitempty"`
}

// Bank represents a financial institution that holds customers and issues cards.
type BankEntitySQL struct {
	ID        uint                `gorm:"primaryKey;autoIncrement"`
	Name      string              `gorm:"size:255"`
	Cuit      string              `gorm:"size:255"`
	Address   string              `gorm:"size:255"`
	Telephone string              `gorm:"size:255"`
	Customers []CustomerEntitySQL `gorm:"many2many:CUSTOMERS_BANKS;"`
	CreatedAt time.Time           `gorm:"autoCreateTime"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime"`
}

func (BankEntitySQL) TableName() string {
	return "BANKS"
}

// ------------ Mappers ------------	//

// Bank a BankModel mapper
func ToBankEntity(bank *models.Bank) *BankEntitySQL {
	return &BankEntitySQL{
		Name:      bank.Name,
		Cuit:      bank.Cuit,
		Address:   bank.Address,
		Telephone: bank.Telephone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BankModel a Bank mapper (si necesitas convertir de nuevo)
func ToBank(bankModel *BankEntitySQL) *models.Bank {
	return &models.Bank{
		Cuit:      bankModel.Cuit,
		Address:   bankModel.Address,
		Telephone: bankModel.Telephone,
		//CustomersIds : []
	}
}

// BankModel NoSQL case overload
func ToBankNonSQL(bank *BankEntityNonSQL) *models.Bank {
	return &models.Bank{
		Cuit:      bank.Cuit,
		Address:   bank.Address,
		Telephone: bank.Telephone,
		//CustomersIds : []
	}
}
