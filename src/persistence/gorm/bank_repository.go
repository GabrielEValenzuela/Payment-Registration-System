package gorm

import (
	"fmt"
	"log"
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

func (r *BankRepositoryGORM) DeleteFinancingPromotion(code string) error {
	var discount entities.FinancingEntity

	// Search for the DiscountEntity by code
	if err := r.db.Where("code = ?", code).First(&discount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("DiscountEntity with code %s not found.", code)
			return err
		}
		log.Printf("Error finding DiscountEntity with code %s: %v", code, err)
		return err
	}

	// Update the IsDeleted field to true
	discount.IsDeleted = true

	// Save changes to the database
	if err := r.db.Save(&discount).Error; err != nil {
		log.Printf("Error updating IsDeleted for DiscountEntity with code %s: %v", code, err)
		return err
	}

	log.Printf("DiscountEntity with code %s was successfully logically deleted.", code)
	return nil
}

func (r *BankRepositoryGORM) DeleteDiscountPromotion(code string) error {
	var financing entities.DiscountEntity

	// Search for the FinancingEntity by code
	if err := r.db.Where("code = ?", code).First(&financing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("FinancingEntity with code %s not found.", code)
			return err
		}
		log.Printf("Error finding FinancingEntity with code %s: %v", code, err)
		return err
	}

	// Update the IsDeleted field to true
	financing.IsDeleted = true

	// Save changes to the database
	if err := r.db.Save(&financing).Error; err != nil {
		log.Printf("Error updating IsDeleted for FinancingEntity with code %s: %v", code, err)
		return err
	}

	log.Printf("FinancingEntity with code %s was successfully logically deleted.", code)
	return nil
}
