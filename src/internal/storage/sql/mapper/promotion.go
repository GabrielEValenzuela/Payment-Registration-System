package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/sql/entities"
)

// Mapper from PromotionModel to Promotion
func ToPromotion(promotionEntity *entities.PromotionEntity) *models.Promotion {
	return &models.Promotion{
		Code:              promotionEntity.Code,
		PromotionTitle:    promotionEntity.PromotionTitle,
		NameStore:         promotionEntity.NameStore,
		CuitStore:         promotionEntity.CuitStore,
		ValidityStartDate: promotionEntity.ValidityStartDate,
		ValidityEndDate:   promotionEntity.ValidityEndDate,
		Comments:          promotionEntity.Comments,
		Bank:              *ToBank(&promotionEntity.Bank), // Use the existing mapper for Bank
	}
}

// Mapper from Promotion to PromotionModel
func ToPromotionEntity(promotion *models.Promotion, bankId uint) *entities.PromotionEntity {
	return &entities.PromotionEntity{
		Code:              promotion.Code,
		PromotionTitle:    promotion.PromotionTitle,
		NameStore:         promotion.NameStore,
		CuitStore:         promotion.CuitStore,
		ValidityStartDate: promotion.ValidityStartDate,
		ValidityEndDate:   promotion.ValidityEndDate,
		Comments:          promotion.Comments,
		Bank:              *ToBankEntity(&promotion.Bank), // Use the existing mapper for BankModel
		BankID:            bankId,                         // Assign the bank's ID
	}
}

// Mapper from FinancingModel to Financing
func ToFinancing(financingEntity *entities.FinancingEntity) *models.Financing {
	return &models.Financing{
		Promotion:      *ToPromotion(&financingEntity.PromotionEntity), // Reuse the PromotionModel mapping
		NumberOfQuotas: financingEntity.NumberOfQuotas,
		Interest:       financingEntity.Interest,
	}
}

// Mapper from Financing to FinancingModel
func ToFinancingEntity(financing *models.Financing, bankId uint) *entities.FinancingEntity {
	return &entities.FinancingEntity{
		PromotionEntity: *ToPromotionEntity(&financing.Promotion, bankId), // Reuse the Promotion mapping
		NumberOfQuotas:  financing.NumberOfQuotas,
		Interest:        financing.Interest,
	}
}

// Mapper from DiscountModel to Discount
func ToDiscount(discountEntity *entities.DiscountEntity) *models.Discount {
	return &models.Discount{
		Promotion:          *ToPromotion(&discountEntity.PromotionEntity), // Reuse the PromotionModel mapping
		DiscountPercentage: discountEntity.DiscountPercentage,
		PriceCap:           discountEntity.DiscountPercentage,
		OnlyCash:           discountEntity.OnlyCash,
	}
}
