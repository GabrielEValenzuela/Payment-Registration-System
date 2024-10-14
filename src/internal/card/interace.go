package card

type Repository interface {
	GetPaymentSummary(cardNumber string, month uint, year uint)
}
