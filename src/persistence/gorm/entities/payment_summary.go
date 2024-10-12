package entities

import "time"

func (PaymentSummaryEntity) TableName() string {
	return "PAYMENT_SUMMARIES"
}

type PaymentSummaryEntity struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement"`
	Code                string     `gorm:"size:255;not null"`
	Month               int        `gorm:"not null"`
	Year                int        `gorm:"not null"`
	FirstExpiration     time.Time  `gorm:"not null"`
	SecondExpiration    time.Time  `gorm:"not null"`
	SurchargePercentage float64    `gorm:"not null"`
	TotalPrice          float64    `gorm:"not null"`
	CardID              uint       `gorm:"not null"`
	Card                CardEntity `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt           time.Time  `gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime"`
}
