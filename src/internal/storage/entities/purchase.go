package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// PurchaseEntity represents the base details of a purchase.
type PurchaseEntityNonSQL struct {
	PaymentVoucher string    `bson:"payment_voucher"`       // Payment voucher code
	Store          string    `bson:"store"`                 // Store name
	CuitStore      string    `bson:"cuit_store"`            // Store CUIT
	Amount         float64   `bson:"amount"`                // Purchase amount
	FinalAmount    float64   `bson:"final_amount"`          // Final amount after adjustments
	CreatedAt      time.Time `bson:"created_at,omitempty"`  // Creation timestamp
	UpdatedAt      time.Time `bson:"updated_at,omitempty"`  // Update timestamp
	CardNumber     string    `bson:"card_number,omitempty"` // Reference to the associated card
}

// PurchaseSinglePaymentEntity represents a single-payment purchase.
type PurchaseSinglePaymentEntityNonSQL struct {
	ID             bson.ObjectID        `bson:"_id,omitempty"`  // MongoDB primary key
	PurchaseEntity PurchaseEntityNonSQL `bson:"purchase"`       // Embedded base purchase details
	StoreDiscount  float64              `bson:"store_discount"` // Discount applied by the store
}

// PurchaseMonthlyPaymentsEntity represents a monthly installment purchase.
type PurchaseMonthlyPaymentsEntityNonSQL struct {
	ID             bson.ObjectID        `bson:"_id,omitempty"`    // MongoDB primary key
	PurchaseEntity PurchaseEntityNonSQL `bson:"purchase"`         // Embedded base purchase details
	Interest       float64              `bson:"interest"`         // Interest rate for the installments
	NumberOfQuotas int                  `bson:"number_of_quotas"` // Number of monthly quotas
	Quotas         []QuotaEntityNonSQL  `bson:"quotas,omitempty"` // Embedded list of quotas
}

type PurchaseEntitySQL struct {
	PaymentVoucher string    `gorm:"size:255;not null"`
	Store          string    `gorm:"size:255;not null"`
	CuitStore      string    `gorm:"size:20;not null"`
	Amount         float64   `gorm:"not null"`
	FinalAmount    float64   `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CardID         uint      `gorm:"index;not null"`
}

type PurchaseSinglePaymentEntitySQL struct {
	ID             uint              `gorm:"primaryKey;autoIncrement"`
	PurchaseEntity PurchaseEntitySQL `gorm:"embedded"`
	StoreDiscount  float64           `gorm:"not null"`
}

type PurchaseMonthlyPaymentsEntitySQL struct {
	ID             uint              `gorm:"primaryKey;autoIncrement"`
	PurchaseEntity PurchaseEntitySQL `gorm:"embedded"`
	Interest       float64           `gorm:"not null"`
	NumberOfQuotas int               `gorm:"not null"`
	Quotas         []QuotaEntitySQL  `gorm:"foreignKey:PurchaseMonthlyPaymentsEntityID"`
}

type PaymentSummaryNoSQL struct {
	Number          string                                `bson:"number"`
	CCV             string                                `bson:"ccv"`
	CardholderName  string                                `bson:"cardholder_name_in_card"`
	BankCuit        string                                `bson:"bank_cuit"`
	CustomerCuit    string                                `bson:"customer_cuit"`
	CreatedAt       time.Time                             `bson:"created_at"`
	UpdatedAt       time.Time                             `bson:"updated_at"`
	SinglePayments  []PurchaseSinglePaymentEntityNonSQL   `bson:"single_payments"`
	MonthlyPayments []PurchaseMonthlyPaymentsEntityNonSQL `bson:"monthly_payments"`
	TotalPrice      float64                               `bson:"total_price"`
}

func (PurchaseSinglePaymentEntitySQL) TableName() string {
	return "PURCHASE_SINGLE_PAYMENTS"
}

func (PurchaseMonthlyPaymentsEntitySQL) TableName() string {
	return "PURCHASE_MONTHLY_PAYMENTS"
}

// ------------ Mappers ------------	//

func ToPurchaseSinglePayment(entity *PurchaseSinglePaymentEntitySQL) *models.PurchaseSinglePayment {
	return &models.PurchaseSinglePayment{
		Purchase:      *toPurchase(&entity.PurchaseEntity),
		StoreDiscount: entity.StoreDiscount,
	}
}

func ToPurchaseSinglePaymentNonSQL(entity *PurchaseSinglePaymentEntityNonSQL) *models.PurchaseSinglePayment {
	return &models.PurchaseSinglePayment{
		Purchase:      *toPurchaseNonSQL(&entity.PurchaseEntity),
		StoreDiscount: entity.StoreDiscount,
	}
}

func ToPurchaseSinglePaymentEntity(model *models.PurchaseSinglePayment) *PurchaseSinglePaymentEntitySQL {
	return &PurchaseSinglePaymentEntitySQL{
		PurchaseEntity: *ToPurchaseEntity(&model.Purchase),
		StoreDiscount:  model.StoreDiscount,
	}
}

func ToPurchaseMonthlyPayments(entity *PurchaseMonthlyPaymentsEntitySQL) *models.PurchaseMonthlyPayment {
	var quotas []models.Quota
	for _, src := range entity.Quotas {
		quotas = append(quotas, *ToQuota(&src))
	}

	return &models.PurchaseMonthlyPayment{
		Purchase:       *toPurchase(&entity.PurchaseEntity),
		NumberOfQuotas: entity.NumberOfQuotas,
		Interest:       entity.Interest,
		Quota:          quotas,
	}
}

func ToPurchaseMonthlyPaymentsNonSQL(entity *PurchaseMonthlyPaymentsEntityNonSQL) *models.PurchaseMonthlyPayment {
	var quotas []models.Quota
	for _, src := range entity.Quotas {
		quotas = append(quotas, *ToQuotaNonSQL(&src))
	}

	return &models.PurchaseMonthlyPayment{
		Purchase:       *toPurchaseNonSQL(&entity.PurchaseEntity),
		NumberOfQuotas: entity.NumberOfQuotas,
		Interest:       entity.Interest,
		Quota:          quotas,
	}
}

func ToPurchaseMonthlyPaymentsEntity(model *models.PurchaseMonthlyPayment) *PurchaseMonthlyPaymentsEntitySQL {
	var quotas []QuotaEntitySQL
	for _, src := range model.Quota {
		quotas = append(quotas, *ToQuotaEntity(&src))
	}
	return &PurchaseMonthlyPaymentsEntitySQL{
		PurchaseEntity: *ToPurchaseEntity(&model.Purchase),
		Interest:       model.Interest,
		NumberOfQuotas: model.NumberOfQuotas,
		Quotas:         quotas,
	}
}

func ToPurchaseEntity(model *models.Purchase) *PurchaseEntitySQL {
	return &PurchaseEntitySQL{
		PaymentVoucher: model.PaymentVoucher,
		Store:          model.Store,
		CuitStore:      model.CuitStore,
		Amount:         model.Amount,
		FinalAmount:    model.FinalAmount,
	}
}

func toPurchase(entity *PurchaseEntitySQL) *models.Purchase {
	return &models.Purchase{
		PaymentVoucher: entity.PaymentVoucher,
		Store:          entity.Store,
		CuitStore:      entity.CuitStore,
		Amount:         entity.Amount,
		FinalAmount:    entity.FinalAmount,
	}
}

func toPurchaseNonSQL(entity *PurchaseEntityNonSQL) *models.Purchase {
	return &models.Purchase{
		PaymentVoucher: entity.PaymentVoucher,
		Store:          entity.Store,
		CuitStore:      entity.CuitStore,
		Amount:         entity.Amount,
		FinalAmount:    entity.FinalAmount,
	}
}

func ConvertPurchaseMonthlyPaymentsList(paymentEntityList *[]PurchaseMonthlyPaymentsEntitySQL) *[]models.PurchaseMonthlyPayment {
	var payments []models.PurchaseMonthlyPayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseMonthlyPayments(&v))
	}
	return &payments
}

func ConvertPurchaseSinglePaymentList(paymentEntityList *[]PurchaseSinglePaymentEntitySQL) *[]models.PurchaseSinglePayment {
	var payments []models.PurchaseSinglePayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseSinglePayment(&v))
	}
	return &payments
}

func ConvertPurchaseSinglePaymentListMongo(paymentEntityList *[]PurchaseSinglePaymentEntityNonSQL) *[]models.PurchaseSinglePayment {
	var payments []models.PurchaseSinglePayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseSinglePaymentNonSQL(&v))
	}
	return &payments
}

func ConvertPurchaseMonthlyPaymentListMongo(paymentEntityList *[]PurchaseMonthlyPaymentsEntityNonSQL) *[]models.PurchaseMonthlyPayment {
	var payments []models.PurchaseMonthlyPayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseMonthlyPaymentsNonSQL(&v))
	}
	return &payments
}
