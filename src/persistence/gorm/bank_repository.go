package gorm

import (
	"fmt"
	"time"

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

func (r *BankRepositoryGORM) ExtendFinancingPromotionValidity(code string, newDate time.Time) error {
	// Find the promotion by ID
	var promotion entities.FinancingEntity
	if err := r.db.First(&promotion, "code = ?", code).Error; err != nil {
		return fmt.Errorf("could not find promotion with code %s: %v", code, err)
	}

	// Update the date
	promotion.ValidityEndDate = newDate

	// Save the changes to the database
	if err := r.db.Save(&promotion).Error; err != nil {
		return fmt.Errorf("could not update promotion dates: %v", err)
	}

	fmt.Printf("Promotion Code %s updated successfully\n", code)
	return nil
}

func (r *BankRepositoryGORM) ExtendDiscountPromotionValidity(code string, newDate time.Time) error {
	// Find the promotion by ID
	var promotion entities.DiscountEntity
	if err := r.db.First(&promotion, "code = ?", code).Error; err != nil {
		return fmt.Errorf("could not find promotion with code %s: %v", code, err)
	}

	// Update the date
	promotion.ValidityEndDate = newDate

	// Save the changes to the database
	if err := r.db.Save(&promotion).Error; err != nil {
		return fmt.Errorf("could not update promotion dates: %v", err)
	}

	fmt.Printf("Promotion Code %s updated successfully\n", code)
	return nil
}
