package gorm

import (
	"log"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/mapper"
	"gorm.io/gorm"
)

type PromotionRepositoryGORM struct {
	db *gorm.DB
}

// NewPromotionRepository crea una nueva instancia de PromotionRepository
func NewPromotionRepository(db *gorm.DB) *PromotionRepositoryGORM {
	return &PromotionRepositoryGORM{db: db}
}

func (r *PromotionRepositoryGORM) GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]promotion.Financing, *[]promotion.Discount, error) {
	var promotionsDiscount []promotion.Discount
	var promotionsFinancing []promotion.Financing

	// Search for DiscountEntity within the date range and matching CUIT
	var discounts []entities.DiscountEntity
	if err := r.db.Where("cuit_store = ? AND is_deleted = ? AND ((validity_start_date >= ? AND validity_end_date <= ?) OR (validity_start_date <= ? AND validity_end_date >= ?))",
		cuit, false, startDate, endDate, startDate, endDate).Find(&discounts).Error; err != nil {
		log.Printf("Error finding DiscountEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}

	// Search for FinancingEntity within the date range and matching CUIT
	var financings []entities.FinancingEntity
	if err := r.db.Where("cuit_store = ? AND is_deleted = ? AND ((validity_start_date >= ? AND validity_end_date <= ?) OR (validity_start_date <= ? AND validity_end_date >= ?))",
		cuit, false, startDate, endDate, startDate, endDate).Find(&financings).Error; err != nil {
		log.Printf("Error finding DiscountEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}

	// Collect all promotions
	for _, discount := range discounts {
		promotionsDiscount = append(promotionsDiscount, *mapper.ToDiscount(&discount))
	}
	for _, financing := range financings {
		promotionsFinancing = append(promotionsFinancing, *mapper.ToFinancing(&financing))
	}

	return &promotionsFinancing, &promotionsDiscount, nil
}
