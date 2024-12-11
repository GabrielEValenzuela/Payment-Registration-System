package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaEntity represents a single installment (quota) for a monthly payment.
type QuotaEntityNonSQL struct {
	ID                              primitive.ObjectID `bson:"_id,omitempty"`                 // MongoDB primary key
	Number                          int                `bson:"number"`                        // Quota number
	Price                           float64            `bson:"price"`                         // Price of the quota
	Month                           string             `bson:"month"`                         // Month of the quota (e.g., "01" for January)
	Year                            string             `bson:"year"`                          // Year of the quota (e.g., "2024")
	PurchaseMonthlyPaymentsEntityID primitive.ObjectID `bson:"purchase_monthly_id,omitempty"` // Reference to the parent PurchaseMonthlyPaymentsEntity
	CreatedAt                       time.Time          `bson:"created_at,omitempty"`          // Creation timestamp
	UpdatedAt                       time.Time          `bson:"updated_at,omitempty"`          // Update timestamp
}

type QuotaEntitySQL struct {
	ID                              uint                             `gorm:"primaryKey;autoIncrement"`
	Number                          int                              `gorm:"not null"`
	Price                           float64                          `gorm:"not null"`
	Month                           string                           `gorm:"size:2;not null"`
	Year                            string                           `gorm:"size:4;not null"`
	PurchaseMonthlyPaymentsEntityID uint                             `gorm:"index;not null"`
	PurchaseMonthlyPaymentsEntity   PurchaseMonthlyPaymentsEntitySQL `gorm:"foreignKey:PurchaseMonthlyPaymentsEntityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt                       time.Time                        `gorm:"autoCreateTime"`
	UpdatedAt                       time.Time                        `gorm:"autoUpdateTime"`
}

func (QuotaEntitySQL) TableName() string {
	return "QUOTAS"
}

// ------------ Mappers ------------	//

func ToQuotaEntity(model *models.Quota) *QuotaEntitySQL {
	return &QuotaEntitySQL{
		Number: model.Number,
		Price:  model.Price,
		Month:  model.Month,
		Year:   model.Year,
	}
}

func ToQuota(entity *QuotaEntitySQL) *models.Quota {
	return &models.Quota{
		Number: entity.Number,
		Price:  entity.Price,
		Month:  entity.Month,
		Year:   entity.Year,
	}
}
