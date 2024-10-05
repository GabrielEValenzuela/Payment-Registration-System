package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/infrastructure/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
)

// Mapper from PromotionModel to Promotion
func ToPromotion(promotionEntity entities.PromotionEntity) promotion.Promotion {
	return promotion.Promotion{
		Code:              promotionEntity.Code,
		PromotionTitle:    promotionEntity.PromotionTitle,
		NameStore:         promotionEntity.NameStore,
		CuitStore:         promotionEntity.CuitStore,
		ValidityStartDate: promotionEntity.ValidityStartDate,
		ValidityEndDate:   promotionEntity.ValidityEndDate,
		Comments:          promotionEntity.Comments,
		Bank:              ToBank(promotionEntity.Bank), // Use the existing mapper for Bank
	}
}

// Mapper from Promotion to PromotionModel
func ToPromotionEntity(promotion promotion.Promotion) entities.PromotionEntity {
	return entities.PromotionEntity{
		Code:              promotion.Code,
		PromotionTitle:    promotion.PromotionTitle,
		NameStore:         promotion.NameStore,
		CuitStore:         promotion.CuitStore,
		ValidityStartDate: promotion.ValidityStartDate,
		ValidityEndDate:   promotion.ValidityEndDate,
		Comments:          promotion.Comments,
		Bank:              ToBankEntity(promotion.Bank), // Use the existing mapper for BankModel
		//BankID:            bankId,                       // Assign the bank's ID
	}
}

// Mapper from FinancingModel to Financing
func ToFinancing(financingEntity entities.FinancingEntity) promotion.Financing {
	return promotion.Financing{
		Promotion:      ToPromotion(financingEntity.PromotionEntity), // Reuse the PromotionModel mapping
		NumberOfQuotas: financingEntity.NumberOfQuotas,
		Interest:       financingEntity.Interest,
	}
}

// Mapper from Financing to FinancingModel
func ToFinancingEntity(financing promotion.Financing) entities.FinancingEntity {
	return entities.FinancingEntity{
		PromotionEntity: ToPromotionEntity(financing.Promotion), // Reuse the Promotion mapping
		NumberOfQuotas:  financing.NumberOfQuotas,
		Interest:        financing.Interest,
	}
}
