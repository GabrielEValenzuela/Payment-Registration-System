package bank

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

type IRepository interface {
	AddFinancingPromotionToBank(promotionFinancing models.Financing) error
	ExtendFinancingPromotionValidity(code string, newDate time.Time) error
	ExtendDiscountPromotionValidity(code string, newDate time.Time) error
	DeleteFinancingPromotion(code string) error
	DeleteDiscountPromotion(code string) error
	GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error)
}
