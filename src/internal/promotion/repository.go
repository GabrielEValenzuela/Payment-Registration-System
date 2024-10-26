package promotion

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

type IRepository interface {
	GetAvailablePromotionsByStoreAndDateRange(cuit string, startDate time.Time, endDate time.Time) (*[]models.Financing, *[]models.Discount, error)
	GetMostUsedPromotion() (*models.Promotion, error)
}
