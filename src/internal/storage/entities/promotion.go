package entities

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ------------------ NoSQL Entities ------------------

// PromotionEntityNonSQL represents a special offer in NoSQL (MongoDB)
type PromotionEntityNonSQL struct {
	Code              string    `bson:"code"`
	PromotionTitle    string    `bson:"promotion_title"`
	NameStore         string    `bson:"name_store"`
	CuitStore         string    `bson:"cuit_store"`
	ValidityStartDate time.Time `bson:"validity_start_date"`
	ValidityEndDate   time.Time `bson:"validity_end_date"`
	Comments          string    `bson:"comments,omitempty"`
}

type FinancingEntityNonSQL struct {
	ID              bson.ObjectID         `bson:"_id,omitempty"`
	PromotionEntity PromotionEntityNonSQL `bson:"promotion_entity"`
	NumberOfQuotas  int                   `bson:"number_of_quotas"`
	Interest        float64               `bson:"interest"`
	IsDeleted       bool                  `bson:"is_deleted"`
	BankID          bson.ObjectID         `bson:"bank_id"`
	CreatedAt       time.Time             `bson:"created_at,omitempty"`
	UpdatedAt       time.Time             `bson:"updated_at,omitempty"`
}

// DiscountEntityNonSQL represents discount promotions in NoSQL
type DiscountEntityNonSQL struct {
	PromotionEntity    PromotionEntityNonSQL `bson:"promotion_entity"`
	DiscountPercentage float64               `bson:"discount_percentage"`
	PriceCap           float64               `bson:"price_cap,omitempty"`
	OnlyCash           bool                  `bson:"only_cash"`
}

// PaymentVoucherCountNonSQL represents voucher usage counts in NoSQL
type PaymentVoucherCountNonSQL struct {
	PaymentVoucher    string `bson:"payment_voucher"`
	TotalRepeticiones int    `bson:"total_repeticiones"`
}

// ------------------ SQL Entities ------------------

// PromotionEntitySQL represents a special offer in SQL (MySQL)
type PromotionEntitySQL struct {
	Code              string        `gorm:"size:255;unique"`
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

// DiscountEntitySQL represents discount promotions in SQL
type DiscountEntitySQL struct {
	PromotionEntitySQL `gorm:"embedded"`
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	DiscountPercentage float64 `gorm:"not null"`
	PriceCap           float64 `gorm:"not null;default:0"`
	OnlyCash           bool    `gorm:"default:false;not null"`
}

// FinancingEntitySQL represents installment-based promotions in SQL
type FinancingEntitySQL struct {
	PromotionEntitySQL `gorm:"embedded"`
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	NumberOfQuotas     int     `gorm:"not null"`
	Interest           float64 `gorm:"not null;default:0"`
}

// PaymentVoucherCountSQL represents voucher usage counts in SQL
type PaymentVoucherCountSQL struct {
	PaymentVoucher    string
	TotalRepeticiones int
}

// ------------------ Table Name Mappings ------------------

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

// ------------------ Mappers ------------------

// SQL -> Models
func ToPromotion(promotionEntity *PromotionEntitySQL) *models.Promotion {
	return &models.Promotion{
		Code:           promotionEntity.Code,
		PromotionTitle: promotionEntity.PromotionTitle,
		NameStore:      promotionEntity.NameStore,
		CuitStore:      promotionEntity.CuitStore,
		ValidityStartDate: models.CustomTime{
			Time: promotionEntity.ValidityStartDate,
		},
		ValidityEndDate: models.CustomTime{
			Time: promotionEntity.ValidityEndDate,
		},
		Comments: promotionEntity.Comments,
		Bank:     *ToBank(&promotionEntity.Bank),
	}
}

func ToPromotionNonSQL(promotionEntity *PromotionEntityNonSQL) *models.Promotion {
	return &models.Promotion{
		Code:           promotionEntity.Code,
		PromotionTitle: promotionEntity.PromotionTitle,
		NameStore:      promotionEntity.NameStore,
		CuitStore:      promotionEntity.CuitStore,
		ValidityStartDate: models.CustomTime{
			Time: promotionEntity.ValidityStartDate,
		},
		ValidityEndDate: models.CustomTime{
			Time: promotionEntity.ValidityEndDate,
		},
		Comments: promotionEntity.Comments,
	}
}

// Models -> SQL
func ToPromotionEntity(promotion *models.Promotion, bankId uint) *PromotionEntitySQL {
	return &PromotionEntitySQL{
		Code:              promotion.Code,
		PromotionTitle:    promotion.PromotionTitle,
		NameStore:         promotion.NameStore,
		CuitStore:         promotion.CuitStore,
		ValidityStartDate: promotion.ValidityStartDate.Time,
		ValidityEndDate:   promotion.ValidityEndDate.Time,
		Comments:          promotion.Comments,
		Bank:              *ToBankEntity(&promotion.Bank),
		BankID:            bankId,
	}
}

// SQL -> Models
func ToFinancing(financingEntity *FinancingEntitySQL) *models.Financing {
	return &models.Financing{
		Promotion:      *ToPromotion(&financingEntity.PromotionEntitySQL),
		NumberOfQuotas: financingEntity.NumberOfQuotas,
		Interest:       financingEntity.Interest,
	}
}

// NoSQL -> Models
func ToFinancingNonSQL(financingEntity *FinancingEntityNonSQL) *models.Financing {
	return &models.Financing{
		Promotion:      *ToPromotionNonSQL(&financingEntity.PromotionEntity),
		NumberOfQuotas: financingEntity.NumberOfQuotas,
		Interest:       financingEntity.Interest,
	}
}

// Models -> SQL
func ToFinancingEntity(financing *models.Financing, bankId uint) *FinancingEntitySQL {
	return &FinancingEntitySQL{
		PromotionEntitySQL: *ToPromotionEntity(&financing.Promotion, bankId),
		NumberOfQuotas:     financing.NumberOfQuotas,
		Interest:           financing.Interest,
	}
}

// SQL -> Models
func ToDiscount(discountEntity *DiscountEntitySQL) *models.Discount {
	return &models.Discount{
		Promotion:          *ToPromotion(&discountEntity.PromotionEntitySQL),
		DiscountPercentage: discountEntity.DiscountPercentage,
		PriceCap:           discountEntity.PriceCap,
		OnlyCash:           discountEntity.OnlyCash,
	}
}

// NoSQL -> Models
func ToDiscountNonSQL(discountEntity *DiscountEntityNonSQL) *models.Discount {
	return &models.Discount{
		Promotion:          *ToPromotionNonSQL(&discountEntity.PromotionEntity),
		DiscountPercentage: discountEntity.DiscountPercentage,
		PriceCap:           discountEntity.PriceCap,
		OnlyCash:           discountEntity.OnlyCash,
	}
}
