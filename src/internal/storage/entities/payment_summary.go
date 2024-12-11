package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentSummaryEntity represents a summary of payments associated with a card.
type PaymentSummaryEntityNonSQL struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`        // MongoDB primary key
	Code                string             `bson:"code"`                 // Unique code for the payment summary
	Month               int                `bson:"month"`                // Payment month
	Year                int                `bson:"year"`                 // Payment year
	FirstExpiration     time.Time          `bson:"first_expiration"`     // First expiration date
	SecondExpiration    time.Time          `bson:"second_expiration"`    // Second expiration date
	SurchargePercentage float64            `bson:"surcharge_percentage"` // Surcharge percentage
	TotalPrice          float64            `bson:"total_price"`          // Total price
	CardID              primitive.ObjectID `bson:"card_id,omitempty"`    // Reference to the associated card
	CreatedAt           time.Time          `bson:"created_at,omitempty"` // Creation timestamp
	UpdatedAt           time.Time          `bson:"updated_at,omitempty"` // Update timestamp
}

type PaymentSummaryEntitySQL struct {
	ID                  uint          `gorm:"primaryKey;autoIncrement"`
	Code                string        `gorm:"size:255;not null"`
	Month               int           `gorm:"not null"`
	Year                int           `gorm:"not null"`
	FirstExpiration     time.Time     `gorm:"not null"`
	SecondExpiration    time.Time     `gorm:"not null"`
	SurchargePercentage float64       `gorm:"not null"`
	TotalPrice          float64       `gorm:"not null"`
	CardID              uint          `gorm:"not null"`
	Card                CardEntitySQL `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt           time.Time     `gorm:"autoCreateTime"`
	UpdatedAt           time.Time     `gorm:"autoUpdateTime"`
}

func (PaymentSummaryEntitySQL) TableName() string {
	return "PAYMENT_SUMMARIES"
}

// ------------ Mappers ------------	//

// Take a model and convert it to a PaymentSummaryEntity for relational storage
func ToPaymentSummaryEntityRelational(paymentSummary *models.PaymentSummary) *PaymentSummaryEntitySQL {
	return &PaymentSummaryEntitySQL{
		Code:                paymentSummary.Code,
		Month:               paymentSummary.Month,
		Year:                paymentSummary.Year,
		FirstExpiration:     paymentSummary.FirstExpiration,
		SecondExpiration:    paymentSummary.SecondExpiration,
		SurchargePercentage: paymentSummary.SurchargePercentage,
		TotalPrice:          paymentSummary.TotalPrice,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

// Take a model and convert it to a PaymentSummaryEntity for non-relational storage
func ToPaymentSummaryEntityNonRelational(paymentSummary *models.PaymentSummary) *PaymentSummaryEntityNonSQL {
	return &PaymentSummaryEntityNonSQL{
		Code:                paymentSummary.Code,
		Month:               paymentSummary.Month,
		Year:                paymentSummary.Year,
		FirstExpiration:     paymentSummary.FirstExpiration,
		SecondExpiration:    paymentSummary.SecondExpiration,
		SurchargePercentage: paymentSummary.SurchargePercentage,
		TotalPrice:          paymentSummary.TotalPrice,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

// PaymentSummaryModel a PaymentSummary mapper (si necesitas convertir de nuevo)
func ToPaymentSummary[T any](paymentSummaryEntity *T) *models.PaymentSummary {
	switch v := any(paymentSummaryEntity).(type) {
	case *PaymentSummaryEntitySQL:
		return &models.PaymentSummary{
			Code:                v.Code,
			Month:               v.Month,
			Year:                v.Year,
			FirstExpiration:     v.FirstExpiration,
			SecondExpiration:    v.SecondExpiration,
			SurchargePercentage: v.SurchargePercentage,
			TotalPrice:          v.TotalPrice,
		}
	case *PaymentSummaryEntityNonSQL:
		return &models.PaymentSummary{
			Code:                v.Code,
			Month:               v.Month,
			Year:                v.Year,
			FirstExpiration:     v.FirstExpiration,
			SecondExpiration:    v.SecondExpiration,
			SurchargePercentage: v.SurchargePercentage,
			TotalPrice:          v.TotalPrice,
		}
	default:
		return nil
	}
}
