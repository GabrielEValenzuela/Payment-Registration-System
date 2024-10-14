package bank

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
)

type Repository interface {
	AddFinancingPromotionToBank(promotionFinancing promotion.Financing) error
	ExtendFinancingPromotionValidity(code string, newDate time.Time) error
	ExtendDiscountPromotionValidity(code string, newDate time.Time) error
}
