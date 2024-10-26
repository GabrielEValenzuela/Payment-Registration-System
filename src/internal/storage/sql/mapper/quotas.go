package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/sql/entities"
)

func ToQuotaEntity(model *models.Quota) *entities.QuotaEntity {
	return &entities.QuotaEntity{
		Number: model.Number,
		Price:  model.Price,
		Month:  model.Month,
		Year:   model.Year,
		//PurchaseMonthlyPaymentsEntityID: purchaseMonthlyPaymentsEntityID,
	}
}

func ToQuota(entity *entities.QuotaEntity) *models.Quota {
	return &models.Quota{
		Number: entity.Number,
		Price:  entity.Price,
		Month:  entity.Month,
		Year:   entity.Year,
	}
}
