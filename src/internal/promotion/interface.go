package promotion

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
)

type Repository interface {
	GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]promotion.Financing, *[]promotion.Discount, error)
	GetMostUsedPromotion() (*promotion.Promotion, error)
}
