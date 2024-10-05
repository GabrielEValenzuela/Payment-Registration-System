package bank

import (
	"context"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
)

type Repository interface {
	AddFinancingPromotionToBank(ctx context.Context, promotionFinancing promotion.Financing) error
}
