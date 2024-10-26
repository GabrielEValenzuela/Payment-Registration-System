package card

import "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"

type IRepository interface {
	GetPaymentSummary(cardNumber string, month int, year int) (*models.PaymentSummary, error)
	GetCardsExpiringInNext30Days(day int, month int, year int) (*[]models.Card, error)
	GetPurchaseMonthly(cuit string, finalAmount int, paymentVoucher string) (*models.PurchaseMonthlyPayment, error)
	GetPurchaseSingle(cuit string, finalAmount int, paymentVoucher string) (*models.PurchaseSinglePayment, error)
	GetTop10CardsByPurchases() (*[]models.Card, error)
}
