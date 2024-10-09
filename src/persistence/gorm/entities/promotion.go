package entities

import (
	"time"
)

func (FinancingEntity) TableName() string {
	return "FINANCING_PROMOTIONS"
}

// Promotion represents a special offer provided by a bank to customers.
// It applies to specific stores and is valid for a certain period of time.
type PromotionEntity struct {
	ID                uint   `gorm:"primaryKey"`
	Code              string `gorm:"size:255"`
	PromotionTitle    string `gorm:"size:255"`
	NameStore         string `gorm:"size:255"`
	CuitStore         string `gorm:"size:255"`
	ValidityStartDate time.Time
	ValidityEndDate   time.Time
	Comments          string     `gorm:"size:255"`
	Bank              BankEntity `gorm:"foreignKey:BankID"`
	BankID            uint       `gorm:"index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Discount represents a type of promotion that applies a percentage discount to purchases.
// Optionally, it may have a maximum discount amount.
type DiscountEntity struct {
	PromotionEntity
	DiscountPercentage float64
	PriceCap           float64
	OnlyCash           bool
}

// Financing represents a promotion that offers installment payment options with specific interest rates.
type FinancingEntity struct {
	PromotionEntity
	NumberOfQuotas int
	Interest       float64
}
