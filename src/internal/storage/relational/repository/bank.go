package relational_repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/entities"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"gorm.io/gorm"
)

type BankRepositoryGORM struct {
	db *gorm.DB
}

func NewBankRelationalRepository(db *gorm.DB) storage.IBankStorage {
	return &BankRepositoryGORM{db: db}
}

// Implementación de la interfaz BankRepository
func (r *BankRepositoryGORM) AddFinancingPromotionToBank(promotionFinancing models.Financing) error {
	var bankEntity entities.BankEntitySQL

	cuit := strings.TrimSpace(promotionFinancing.Bank.Cuit)
	logger.Info("Searching for bank with 'cuit' %s", cuit)
	if err := r.db.First(&bankEntity, "cuit = ?", cuit).Error; err != nil {
		return fmt.Errorf("could not find bank with 'cuit' %s: %v", promotionFinancing.Bank.Cuit, err)
	}

	// Lógica para agregar la promoción al banco
	return r.db.Create(entities.ToFinancingEntity(&promotionFinancing, bankEntity.ID)).Error
}

func (r *BankRepositoryGORM) ExtendFinancingPromotionValidity(code string, newDate time.Time) error {
	// Find the promotion by ID
	var promotion entities.FinancingEntitySQL
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
	var promotion entities.DiscountEntitySQL
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
	var discount entities.FinancingEntitySQL

	// Search for the DiscountEntity by code
	if err := r.db.Where("code = ?", code).First(&discount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("DiscountEntity with code %s not found.", code)
			return err
		}
		logger.Info("Error finding DiscountEntity with code %s: %v", code, err)
		return err
	}

	// Update the IsDeleted field to true
	discount.IsDeleted = true

	// Save changes to the database
	if err := r.db.Save(&discount).Error; err != nil {
		logger.Info("Error updating IsDeleted for DiscountEntity with code %s: %v", code, err)
		return err
	}

	logger.Info("DiscountEntity with code %s was successfully logically deleted.", code)
	return nil
}

func (r *BankRepositoryGORM) DeleteDiscountPromotion(code string) error {
	var financing entities.DiscountEntitySQL

	// Search for the FinancingEntity by code
	if err := r.db.Where("code = ?", code).First(&financing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("FinancingEntity with code %s not found.", code)
			return err
		}
		logger.Info("Error finding FinancingEntity with code %s: %v", code, err)
		return err
	}

	// Update the IsDeleted field to true
	financing.IsDeleted = true

	// Save changes to the database
	if err := r.db.Save(&financing).Error; err != nil {
		logger.Info("Error updating IsDeleted for FinancingEntity with code %s: %v", code, err)
		return err
	}

	logger.Info("FinancingEntity with code %s was successfully logically deleted.", code)
	return nil
}

func (r *BankRepositoryGORM) GetBankCustomerCounts() ([]models.BankCustomerCountDTO, error) {
	var results []models.BankCustomerCountDTO

	r.db.Raw(`
		SELECT 
			b.cuit AS bank_cuit,
			b.name AS bank_name,
			COUNT(cb.customer_entity_sql_id) AS customer_count
		FROM 
			BANKS b
		LEFT JOIN 
			customers_banks cb ON b.id = cb.bank_entity_sql_id
		GROUP BY 
			b.id, b.name
	`).Scan(&results)

	return results, nil
}
