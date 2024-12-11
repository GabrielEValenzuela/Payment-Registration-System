package relational

import (
	"errors"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"gorm.io/gorm"
)

type PromotionRepositoryGORM struct {
	db *gorm.DB
}

// NewPromotionRepository crea una nueva instancia de PromotionRepository
func NewPromotionRepository(db *gorm.DB) *PromotionRepositoryGORM {
	return &PromotionRepositoryGORM{db: db}
}

func (r *PromotionRepositoryGORM) GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]models.Financing, *[]models.Discount, error) {
	var promotionsDiscount []models.Discount
	var promotionsFinancing []models.Financing

	// Search for DiscountEntity within the date range and matching CUIT
	var discounts []entities.DiscountEntitySQL
	if err := r.db.Where("cuit_store = ? AND is_deleted = ? AND ((validity_start_date >= ? AND validity_end_date <= ?) OR (validity_start_date <= ? AND validity_end_date >= ?))",
		cuit, false, startDate, endDate, startDate, endDate).Find(&discounts).Error; err != nil {
		logger.Info("Error finding DiscountEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}

	// Search for FinancingEntity within the date range and matching CUIT
	var financings []entities.FinancingEntitySQL
	if err := r.db.Where("cuit_store = ? AND is_deleted = ? AND ((validity_start_date >= ? AND validity_end_date <= ?) OR (validity_start_date <= ? AND validity_end_date >= ?))",
		cuit, false, startDate, endDate, startDate, endDate).Find(&financings).Error; err != nil {
		logger.Info("Error finding DiscountEntity with CUIT %s between %v and %v: %v", cuit, startDate, endDate, err)
		return nil, nil, err
	}

	// Collect all promotions
	for _, discount := range discounts {
		promotionsDiscount = append(promotionsDiscount, *entities.ToDiscount(&discount))
	}
	for _, financing := range financings {
		promotionsFinancing = append(promotionsFinancing, *entities.ToFinancing(&financing))
	}

	return &promotionsFinancing, &promotionsDiscount, nil
}

// GetMostUsedDiscountPromotion retrieves the most used discount promotion based on its usage in single and monthly payments.
func (r *PromotionRepositoryGORM) GetMostUsedPromotion() (interface{}, error) {
	var result entities.PaymentVoucherCountSQL

	query := `
		SELECT
			payment_voucher,
			COUNT(*) AS total_repeticiones
		FROM
			(
			SELECT
				month.payment_voucher
			FROM
				PURCHASES_MONTHLY_PAYMENTS month
			UNION ALL
			SELECT
				single.payment_voucher
			FROM
				PURCHASES_SINGLE_PAYMENTS single
			) as payment_voucher
		GROUP BY
			payment_voucher
		ORDER BY
			total_repeticiones DESC
		LIMIT 1;
		`

	r.db.Raw(query).Scan(&result)

	promotion, err := findPromotionByCode(r.db, result.PaymentVoucher)
	if err != nil {
		return nil, err
	}

	return promotion, nil
}

func findPromotionByCode(db *gorm.DB, code string) (interface{}, error) {
	var discountPromo entities.DiscountEntitySQL
	var financingPromo entities.FinancingEntitySQL

	if err := db.Where("code = ?", code).First(&financingPromo).Error; err == nil {
		return financingPromo, nil
	}

	if err := db.Where("code = ?", code).First(&discountPromo).Error; err == nil {
		return discountPromo, nil
	}

	return nil, errors.New("promotion not found with the provided code")
}
