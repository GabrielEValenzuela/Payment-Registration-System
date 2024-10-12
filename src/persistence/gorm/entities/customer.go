package entities

import (
	"time"
)

func (CustomerEntity) TableName() string {
	return "CUSTOMERS"
}

type CustomerEntity struct {
	ID           uint         `gorm:"primaryKey;autoIncrement"`
	CompleteName string       `gorm:"size:255;not null"`
	Dni          string       `gorm:"size:20;unique;not null"`
	Cuit         string       `gorm:"size:20;unique;not null"`
	Address      string       `gorm:"size:255"`
	Telephone    string       `gorm:"size:50"`
	EntryDate    time.Time    `gorm:"not null"`
	Banks        []BankEntity `gorm:"many2many:CUSTOMERS_BANKS"`
	Cards        []CardEntity `gorm:"foreignKey:CustomerID"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime"`
}
