package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/sql/entities"
)

func ToCard(entity *entities.CardEntity) *models.Card {
	return &models.Card{
		Number:                  entity.Number,
		Ccv:                     entity.Ccv,
		CardholderNameInCard:    entity.CardholderNameInCard,
		Since:                   entity.Since,
		ExpirationDate:          entity.ExpirationDate,
		Bank:                    *ToBank(&entity.Bank),
		PurchaseMonthlyPayments: *ConvertPurchaseMonthlyPaymentsList(&entity.PurchaseMonthlyPayments),
		PurchaseSinglePayments:  *ConvertPurchaseSinglePaymentList(&entity.PurchaseSinglePayments),
	}
}

func ToCardEntity(model *models.Card) *entities.CardEntity {
	return &entities.CardEntity{
		Number:               model.Number,
		Ccv:                  model.Ccv,
		CardholderNameInCard: model.CardholderNameInCard,
		Since:                model.Since,
		ExpirationDate:       model.ExpirationDate,
		Bank:                 *ToBankEntity(&model.Bank),
	}
}
