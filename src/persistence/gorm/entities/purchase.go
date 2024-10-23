package entities

import (
	"time"
)

func (PurchaseSinglePaymentEntity) TableName() string {
	return "PURCHASES_SINGLE_PAYMENTS"
}

func (PurchaseMonthlyPaymentsEntity) TableName() string {
	return "PURCHASES_MONTHLY_PAYMENTS"
}

type PurchaseEntity struct {
	PaymentVoucher string    `gorm:"size:255;not null"`
	Store          string    `gorm:"size:255;not null"`
	CuitStore      string    `gorm:"size:20;not null"`
	Amount         float64   `gorm:"not null"`
	FinalAmount    float64   `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CardID         uint      `gorm:"index;not null"`
}

type PurchaseSinglePaymentEntity struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	PurchaseEntity PurchaseEntity `gorm:"embedded"`
	StoreDiscount  float64        `gorm:"not null"`
}

type PurchaseMonthlyPaymentsEntity struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	PurchaseEntity PurchaseEntity `gorm:"embedded"`
	Interest       float64        `gorm:"not null"`
	NumberOfQuotas int            `gorm:"not null"`
	Quotas         []QuotaEntity  `gorm:"foreignKey:PurchaseMonthlyPaymentsEntityID"`
}
