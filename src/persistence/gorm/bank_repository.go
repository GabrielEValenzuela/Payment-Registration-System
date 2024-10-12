package gorm

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/promotion"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/mapper"
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
	var bankEntity entities.BankEntity

	if err := r.db.First(&bankEntity, "cuit = ?", promotionFinancing.Bank.Cuit).Error; err != nil {
		panic(err)
	}

	// Lógica para agregar la promoción al banco
	return r.db.Create(mapper.ToFinancingEntity(promotionFinancing, bankEntity.ID)).Error
}
