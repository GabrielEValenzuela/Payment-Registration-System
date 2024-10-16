package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/payment_summary"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
)

func ToPaymentSummary(entity *entities.PaymentSummaryEntity) *payment_summary.PaymentSummary {
	return &payment_summary.PaymentSummary{
		Code:                entity.Code,
		Month:               entity.Month,
		Year:                entity.Year,
		FirstExpiration:     entity.FirstExpiration,
		SecondExpiration:    entity.SecondExpiration,
		SurchargePercentage: entity.SurchargePercentage,
		TotalPrice:          entity.TotalPrice,
		Card:                *ToCard(&entity.Card),
		MonthlyPayments:     *ConvertPurchaseMonthlyPaymentsList(&entity.Card.PurchaseMonthlyPayments),
		SinglePayments:      *ConvertPurchaseSinglePaymentList(&entity.Card.PurchaseSinglePayments),
	}
}

func PaymentSummaryEntity(model *payment_summary.PaymentSummary, cardId uint) *entities.PaymentSummaryEntity {
	return &entities.PaymentSummaryEntity{
		Code:                model.Code,
		Month:               model.Month,
		Year:                model.Year,
		FirstExpiration:     model.FirstExpiration,
		SecondExpiration:    model.SecondExpiration,
		SurchargePercentage: model.SurchargePercentage,
		TotalPrice:          model.TotalPrice,
		Card:                *ToCardEntity(&model.Card),
		CardID:              cardId,
	}
}
