package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/quota"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
)

func ToQuotaEntity(model *quota.Quota) *entities.QuotaEntity {
	return &entities.QuotaEntity{
		Number: model.Number,
		Price:  model.Price,
		Month:  model.Month,
		Year:   model.Year,
		//PurchaseMonthlyPaymentsEntityID: purchaseMonthlyPaymentsEntityID,
	}
}

func ToQuota(entity *entities.QuotaEntity) *quota.Quota {
	return &quota.Quota{
		Number: entity.Number,
		Price:  entity.Price,
		Month:  entity.Month,
		Year:   entity.Year,
	}
}
