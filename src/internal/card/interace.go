package card

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/card"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/payment_summary"
)

type Repository interface {
	GetPaymentSummary(cardNumber string, month int, year int) (payment_summary.PaymentSummary, error)
	GetCardsExpiringInNext30Days(day int, month int, year int) ([]card.Card, error)
}
