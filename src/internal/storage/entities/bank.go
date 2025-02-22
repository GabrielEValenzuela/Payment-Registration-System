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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
 * BankEntityNonSQL
 * ----------------
 * Represents a financial institution that holds customers and issues cards.
 * The entity is implemented using MongoDB.
 * The entity has a one-to-many relationship with the Customer entity.
 *
 * Fields:
 * - ID (primitive.ObjectID): Unique identifier for the bank
 * - Name (string): Name of the bank
 * - Cuit (string): Unique tax identification code for the bank
 * - Address (string): Physical address of the bank
 * - Telephone (string): Contact number for the bank
 * - Customers ([]CustomerEntityNonSQL): List of associated customers
 * - CreatedAt (time.Time): Timestamp of creation
 * - UpdatedAt (time.Time): Timestamp of last update
 */
type BankEntityNonSQL struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty"` // Use ObjectId as the primary key
	Name      string                 `bson:"name"`
	Cuit      string                 `bson:"cuit;unique"`
	Address   string                 `bson:"address"`
	Telephone string                 `bson:"telephone"`
	Customers []CustomerEntityNonSQL `bson:"customers,omitempty"`
	CreatedAt time.Time              `bson:"created_at,omitempty"`
	UpdatedAt time.Time              `bson:"updated_at,omitempty"`
}

/*
* BankEntitySQL
* -------------
*
* Represents a financial institution that holds customers and issues cards.
* The entity is implemented using GORM.
* The entity has a many-to-many relationship with the Customer entity.
*
* Fields:
* - ID (uint): Unique identifier for the bank
* - Name (string): Name of the bank
* - Cuit (string): Unique tax identification code for the bank
* - Address (string): Physical address of the bank
* - Telephone (string): Contact number for the bank
* - Customers ([]CustomerEntitySQL): List of associated customers
* - CreatedAt (time.Time): Timestamp of creation
* - UpdatedAt (time.Time): Timestamp of last update
 */
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
