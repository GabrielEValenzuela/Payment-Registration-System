package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
)

func ToCard(entity *entities.CardEntity) *card.Card {
	return &card.Card{
		Number:               entity.Number,
		Ccv:                  entity.Ccv,
		CardholderNameInCard: entity.CardholderNameInCard,
		Since:                entity.Since,
		ExpirationDate:       entity.ExpirationDate,
		Bank:                 *ToBank(&entity.Bank),
		//Purchases: ,
	}
}

func ToCardEntity(model *card.Card) *entities.CardEntity {
	return &entities.CardEntity{
		Number:               model.Number,
		Ccv:                  model.Ccv,
		CardholderNameInCard: model.CardholderNameInCard,
		Since:                model.Since,
		ExpirationDate:       model.ExpirationDate,
		Bank:                 *ToBankEntity(&model.Bank),
	}
}
