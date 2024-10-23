package entities

import (
	"time"
)

func (FinancingEntity) TableName() string {
	return "FINANCING_PROMOTIONS"
}

func (DiscountEntity) TableName() string {
	return "DISCOUNT_PROMOTIONS"
}

// Promotion represents a special offer provided by a bank to customers.
// It applies to specific stores and is valid for a certain period of time.
type PromotionEntity struct {
	Code              string     `gorm:"size:255"`
	PromotionTitle    string     `gorm:"size:255"`
	NameStore         string     `gorm:"size:255"`
	CuitStore         string     `gorm:"size:255"`
	ValidityStartDate time.Time  `gorm:"not null"`
	ValidityEndDate   time.Time  `gorm:"not null"`
	Comments          string     `gorm:"size:255"`
	Bank              BankEntity `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BankID            uint       `gorm:"index"`
	CreatedAt         time.Time  `gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime"`
	IsDeleted         bool       `gorm:"default:false;not null"`
	PurchaseCount     int        `gorm:"-"`
}

// Discount represents a type of promotion that applies a percentage discount to purchases.
// Optionally, it may have a maximum discount amount.
type DiscountEntity struct {
	PromotionEntity    `gorm:"embedded"`
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	DiscountPercentage float64 `gorm:"not null"`
	PriceCap           float64 `gorm:"not null;default:0"`
	OnlyCash           bool    `gorm:"default:false;not null"`
}

// Financing represents a promotion that offers installment payment options with specific interest rates.
type FinancingEntity struct {
	PromotionEntity `gorm:"embedded"`
	ID              uint    `gorm:"primaryKey;autoIncrement"`
	NumberOfQuotas  int     `gorm:"not null"`
	Interest        float64 `gorm:"not null;default:0"`
}

type PaymentVoucherCount struct {
	PaymentVoucher    string
	TotalRepeticiones int
}
