package entities

import "time"

func (QuotaEntity) TableName() string {
	return "QUOTAS"
}

type QuotaEntity struct {
	ID                              uint                          `gorm:"primaryKey;autoIncrement"`
	Number                          int                           `gorm:"not null"`
	Price                           float64                       `gorm:"not null"`
	Month                           string                        `gorm:"size:2;not null"`
	Year                            string                        `gorm:"size:4;not null"`
	PurchaseMonthlyPaymentsEntityID uint                          `gorm:"index;not null"`
	PurchaseMonthlyPaymentsEntity   PurchaseMonthlyPaymentsEntity `gorm:"foreignKey:PurchaseMonthlyPaymentsEntityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt                       time.Time                     `gorm:"autoCreateTime"`
	UpdatedAt                       time.Time                     `gorm:"autoUpdateTime"`
}
