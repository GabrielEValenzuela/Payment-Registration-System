package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CardEntity represents a credit or debit card issued by a bank.
type CardEntityNonSQL struct {
	ID                      bson.ObjectID                         `bson:"_id,omitempty"` // MongoDB primary key
	Number                  string                                `bson:"number"`        // Card number (16 digits)
	Ccv                     string                                `bson:"ccv"`           // Card verification code (3 digits)
	CardholderNameInCard    string                                `bson:"cardholder_name_in_card"`
	Since                   time.Time                             `bson:"since"` // When the card was issued
	ExpirationDate          time.Time                             `bson:"expiration_date"`
	BankCuit                string                                `bson:"bank_cuit,omitempty"`     // Reference to the bank (if using references)
	CustomerCuit            string                                `bson:"customer_cuit,omitempty"` // Reference to the customer
	PurchaseSinglePayments  []PurchaseSinglePaymentEntityNonSQL   `bson:"purchase_single_payments,omitempty"`
	PurchaseMonthlyPayments []PurchaseMonthlyPaymentsEntityNonSQL `bson:"purchase_monthly_payments,omitempty"`
	CreatedAt               time.Time                             `bson:"created_at,omitempty"` // Creation timestamp
	UpdatedAt               time.Time                             `bson:"updated_at,omitempty"` // Update timestamp
}

type PaymentSummaryResult struct {
	Card                    CardEntityNonSQL `bson:",inline"`
	PurchaseSinglePayments  []bson.Raw       `bson:"purchase_single_payments"`
	PurchaseMonthlyPayments []bson.Raw       `bson:"purchase_monthly_payments"`
	TotalAmount             float64          `bson:"total_amount"`
}

type CardEntitySQL struct {
	ID                      uint                               `gorm:"primaryKey;autoIncrement"`
	Number                  string                             `gorm:"size:16;not null"`
	Ccv                     string                             `gorm:"size:3;not null"`
	CardholderNameInCard    string                             `gorm:"size:255;not null"`
	Since                   time.Time                          `gorm:"not null"`
	ExpirationDate          time.Time                          `gorm:"not null"`
	Bank                    BankEntitySQL                      `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BankID                  uint                               `gorm:"index"`
	CustomerID              uint                               `gorm:"index"`
	PurchaseSinglePayments  []PurchaseSinglePaymentEntitySQL   `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PurchaseMonthlyPayments []PurchaseMonthlyPaymentsEntitySQL `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt               time.Time                          `gorm:"autoCreateTime"`
	UpdatedAt               time.Time                          `gorm:"autoUpdateTime"`
}

func (CardEntitySQL) TableName() string {
	return "CARDS"
}

// ------------ Mappers ------------	//

// Take a model and convert it to a BankEntity for relational storage
func ToCardEntityRelational(card *models.Card) *CardEntitySQL {
	return &CardEntitySQL{
		Number:               card.Number,
		Ccv:                  card.Ccv,
		CardholderNameInCard: card.CardholderNameInCard,
		Since:                card.Since,
		ExpirationDate:       card.ExpirationDate,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// Take a model and convert it to a BankEntity for non-relational storage
func ToCardEntityNonRelational(card *models.Card) *CardEntityNonSQL {
	return &CardEntityNonSQL{
		Number:               card.Number,
		Ccv:                  card.Ccv,
		CardholderNameInCard: card.CardholderNameInCard,
		Since:                card.Since,
		ExpirationDate:       card.ExpirationDate,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// CardModel a Card mapper (si necesitas convertir de nuevo)
func ToCard[T any](cardEntity *T) *models.Card {
	switch v := any(cardEntity).(type) {
	case *CardEntitySQL:
		return &models.Card{
			Number:                  v.Number,
			Ccv:                     v.Ccv,
			CardholderNameInCard:    v.CardholderNameInCard,
			Since:                   v.Since,
			ExpirationDate:          v.ExpirationDate,
			Bank:                    *ToBank(&v.Bank),
			PurchaseMonthlyPayments: *ConvertPurchaseMonthlyPaymentsList(&v.PurchaseMonthlyPayments),
			PurchaseSinglePayments:  *ConvertPurchaseSinglePaymentList(&v.PurchaseSinglePayments),
		}
	case *CardEntityNonSQL:
		return &models.Card{
			Number:                  v.Number,
			Ccv:                     v.Ccv,
			CardholderNameInCard:    v.CardholderNameInCard,
			Since:                   v.Since,
			ExpirationDate:          v.ExpirationDate,
			PurchaseMonthlyPayments: *ConvertPurchaseMonthlyPaymentListMongo(&v.PurchaseMonthlyPayments),
			PurchaseSinglePayments:  *ConvertPurchaseSinglePaymentListMongo(&v.PurchaseSinglePayments),
		}
	}
	return nil
}
