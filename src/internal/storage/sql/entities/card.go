package entities

import (
	"time"
)

func (CardEntity) TableName() string {
	return "CARDS"
}

type CardEntity struct {
	ID                      uint                            `gorm:"primaryKey;autoIncrement"`
	Number                  string                          `gorm:"size:16;not null"`
	Ccv                     string                          `gorm:"size:3;not null"`
	CardholderNameInCard    string                          `gorm:"size:255;not null"`
	Since                   time.Time                       `gorm:"not null"`
	ExpirationDate          time.Time                       `gorm:"not null"`
	Bank                    BankEntity                      `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BankID                  uint                            `gorm:"index"`
	CustomerID              uint                            `gorm:"index"`
	PurchaseSinglePayments  []PurchaseSinglePaymentEntity   `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PurchaseMonthlyPayments []PurchaseMonthlyPaymentsEntity `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt               time.Time                       `gorm:"autoCreateTime"`
	UpdatedAt               time.Time                       `gorm:"autoUpdateTime"`
}
