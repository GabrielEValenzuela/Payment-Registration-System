package gorm

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/infrastructure/persistence/gorm/mapper"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
	"gorm.io/gorm"
)

type BankRepositoryGORM struct {
	db *gorm.DB
}

// NewBankRepository crea una nueva instancia de BankRepository
func NewBankRepository(db *gorm.DB) *BankRepositoryGORM {
	return &BankRepositoryGORM{db: db}
}

// Implementación de la interfaz BankRepository
func (r *BankRepositoryGORM) AddFinancingPromotionToBank(promotionFinancing *promotion.Financing) error {
	// Lógica para agregar la promoción al banco
	return r.db.Create(mapper.ToFinancingEntity(promotionFinancing)).Error
}
