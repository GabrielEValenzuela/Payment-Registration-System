package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PromotionEntity represents a special offer provided by a bank to customers.
// It applies to specific stores and is valid for a certain period of time.
type PromotionEntityNonSQL struct {
	Code              string             `bson:"code"`                     // Unique code for the promotion
	PromotionTitle    string             `bson:"promotion_title"`          // Title of the promotion
	NameStore         string             `bson:"name_store"`               // Store name
	CuitStore         string             `bson:"cuit_store"`               // Store CUIT
	ValidityStartDate time.Time          `bson:"validity_start_date"`      // Start date of validity
	ValidityEndDate   time.Time          `bson:"validity_end_date"`        // End date of validity
	Comments          string             `bson:"comments,omitempty"`       // Optional comments
	BankID            primitive.ObjectID `bson:"bank_id,omitempty"`        // Reference to associated bank
	CreatedAt         time.Time          `bson:"created_at,omitempty"`     // Creation timestamp
	UpdatedAt         time.Time          `bson:"updated_at,omitempty"`     // Update timestamp
	IsDeleted         bool               `bson:"is_deleted"`               // Soft delete flag
	PurchaseCount     int                `bson:"purchase_count,omitempty"` // Count of purchases (not persisted)
}

// DiscountEntity represents a type of promotion that applies a percentage discount to purchases.
// Optionally, it may have a maximum discount amount.
type DiscountEntityNonSQL struct {
	ID                 primitive.ObjectID    `bson:"_id,omitempty"`       // MongoDB primary key
	PromotionEntity    PromotionEntityNonSQL `bson:"promotion_entity"`    // Embedded promotion details
	DiscountPercentage float64               `bson:"discount_percentage"` // Discount percentage
	PriceCap           float64               `bson:"price_cap,omitempty"` // Optional maximum discount amount
	OnlyCash           bool                  `bson:"only_cash"`           // Only cash flag
}

// FinancingEntity represents a promotion that offers installment payment options with specific interest rates.
type FinancingEntityNonSQL struct {
	ID              primitive.ObjectID    `bson:"_id,omitempty"`    // MongoDB primary key
	PromotionEntity PromotionEntityNonSQL `bson:"promotion_entity"` // Embedded promotion details
	NumberOfQuotas  int                   `bson:"number_of_quotas"` // Number of installment payments
	Interest        float64               `bson:"interest"`         // Interest rate
}

// PaymentVoucherCount represents a summary of payment voucher usage counts.
type PaymentVoucherCountNonSQL struct {
	PaymentVoucher    string `bson:"payment_voucher"`    // Payment voucher identifier
	TotalRepeticiones int    `bson:"total_repeticiones"` // Total number of repetitions
}

// Promotion represents a special offer provided by a bank to customers.
// It applies to specific stores and is valid for a certain period of time.
type PromotionEntitySQL struct {
	Code              string        `gorm:"size:255"`
	PromotionTitle    string        `gorm:"size:255"`
	NameStore         string        `gorm:"size:255"`
	CuitStore         string        `gorm:"size:255"`
	ValidityStartDate time.Time     `gorm:"not null"`
	ValidityEndDate   time.Time     `gorm:"not null"`
	Comments          string        `gorm:"size:255"`
	Bank              BankEntitySQL `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BankID            uint          `gorm:"index"`
	CreatedAt         time.Time     `gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime"`
	IsDeleted         bool          `gorm:"default:false;not null"`
	PurchaseCount     int           `gorm:"-"`
}

// Discount represents a type of promotion that applies a percentage discount to purchases.
// Optionally, it may have a maximum discount amount.
type DiscountEntitySQL struct {
	PromotionEntitySQL `gorm:"embedded"`
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	DiscountPercentage float64 `gorm:"not null"`
	PriceCap           float64 `gorm:"not null;default:0"`
	OnlyCash           bool    `gorm:"default:false;not null"`
}

// Financing represents a promotion that offers installment payment options with specific interest rates.
type FinancingEntitySQL struct {
	PromotionEntitySQL `gorm:"embedded"`
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	NumberOfQuotas     int     `gorm:"not null"`
	Interest           float64 `gorm:"not null;default:0"`
}

type PaymentVoucherCountSQL struct {
	PaymentVoucher    string
	TotalRepeticiones int
}

func (PromotionEntitySQL) TableName() string {
	return "PROMOTIONS"
}

func (DiscountEntitySQL) TableName() string {
	return "DISCOUNTS"
}

func (FinancingEntitySQL) TableName() string {
	return "FINANCINGS"
}

func (PaymentVoucherCountSQL) TableName() string {
	return "PAYMENT_VOUCHER_COUNTS"
}

// ------------ Mappers ------------	//

// Mapper from PromotionModel to Promotion
func ToPromotion(promotionEntity *PromotionEntitySQL) *models.Promotion {
	return &models.Promotion{
		Code:              promotionEntity.Code,
		PromotionTitle:    promotionEntity.PromotionTitle,
		NameStore:         promotionEntity.NameStore,
		CuitStore:         promotionEntity.CuitStore,
		ValidityStartDate: promotionEntity.ValidityStartDate,
		ValidityEndDate:   promotionEntity.ValidityEndDate,
		Comments:          promotionEntity.Comments,
		Bank:              *ToBank(&promotionEntity.Bank), // Use the existing mapper for Bank
	}
}

func ToPromotionNonSQL(promotionEntity *PromotionEntityNonSQL) *models.Promotion {
	return &models.Promotion{
		Code:              promotionEntity.Code,
		PromotionTitle:    promotionEntity.PromotionTitle,
		NameStore:         promotionEntity.NameStore,
		CuitStore:         promotionEntity.CuitStore,
		ValidityStartDate: promotionEntity.ValidityStartDate,
		ValidityEndDate:   promotionEntity.ValidityEndDate,
		Comments:          promotionEntity.Comments,
	}
}

// Mapper from Promotion to PromotionModel
func ToPromotionEntity(promotion *models.Promotion, bankId uint) *PromotionEntitySQL {
	return &PromotionEntitySQL{
		Code:              promotion.Code,
		PromotionTitle:    promotion.PromotionTitle,
		NameStore:         promotion.NameStore,
		CuitStore:         promotion.CuitStore,
		ValidityStartDate: promotion.ValidityStartDate,
		ValidityEndDate:   promotion.ValidityEndDate,
		Comments:          promotion.Comments,
		Bank:              *ToBankEntity(&promotion.Bank), // Use the existing mapper for BankModel
		BankID:            bankId,                         // Assign the bank's ID
	}
}

// Mapper from FinancingModel to Financing
func ToFinancing(financingEntity *FinancingEntitySQL) *models.Financing {
	return &models.Financing{
		Promotion:      *ToPromotion(&financingEntity.PromotionEntitySQL), // Reuse the PromotionModel mapping
		NumberOfQuotas: financingEntity.NumberOfQuotas,
		Interest:       financingEntity.Interest,
	}
}

func ToFinancingNonSQL(financingEntity *FinancingEntityNonSQL) *models.Financing {
	return &models.Financing{
		Promotion:      *ToPromotionNonSQL(&financingEntity.PromotionEntity),
		NumberOfQuotas: financingEntity.NumberOfQuotas,
		Interest:       financingEntity.Interest,
	}
}

// Mapper from Financing to FinancingModel
func ToFinancingEntity(financing *models.Financing, bankId uint) *FinancingEntitySQL {
	return &FinancingEntitySQL{
		PromotionEntitySQL: *ToPromotionEntity(&financing.Promotion, bankId), // Reuse the Promotion mapping
		NumberOfQuotas:     financing.NumberOfQuotas,
		Interest:           financing.Interest,
	}
}

// Mapper from DiscountModel to Discount
func ToDiscount(discountEntity *DiscountEntitySQL) *models.Discount {
	return &models.Discount{
		Promotion:          *ToPromotion(&discountEntity.PromotionEntitySQL), // Reuse the PromotionModel mapping
		DiscountPercentage: discountEntity.DiscountPercentage,
		PriceCap:           discountEntity.DiscountPercentage,
		OnlyCash:           discountEntity.OnlyCash,
	}
}

func ToDiscountNonSQL(discountEntity *DiscountEntityNonSQL) *models.Discount {
	return &models.Discount{
		Promotion:          *ToPromotionNonSQL(&discountEntity.PromotionEntity),
		DiscountPercentage: discountEntity.DiscountPercentage,
		PriceCap:           discountEntity.PriceCap,
		OnlyCash:           discountEntity.OnlyCash,
	}
}
