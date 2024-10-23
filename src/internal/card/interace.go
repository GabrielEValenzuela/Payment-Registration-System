package card

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/payment_summary"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/purchase"
)

type Repository interface {
	GetPaymentSummary(cardNumber string, month int, year int) (*payment_summary.PaymentSummary, error)
	GetCardsExpiringInNext30Days(day int, month int, year int) (*[]card.Card, error)
	GetPurchaseMonthly(cuit string, finalAmount int, paymentVoucher string) (*purchase.PurchaseMonthlyPayment, error)
	GetPurchaseSingle(cuit string, finalAmount int, paymentVoucher string) (*purchase.PurchaseSinglePayment, error)
	GetTop10CardsByPurchases() (*[]card.Card, error)
}
