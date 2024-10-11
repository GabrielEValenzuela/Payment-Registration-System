package mapper

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/bank"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
)

// Bank a BankModel mapper
func ToBankEntity(bank *bank.Bank) *entities.BankEntity {
	return &entities.BankEntity{
		Name:      bank.Name,
		Cuit:      bank.Cuit,
		Address:   bank.Address,
		Telephone: bank.Telephone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BankModel a Bank mapper (si necesitas convertir de nuevo)
func ToBank(bankModel *entities.BankEntity) *bank.Bank {
	return &bank.Bank{
		Cuit:      bankModel.Cuit,
		Address:   bankModel.Address,
		Telephone: bankModel.Telephone,
		//CustomersIds : []
	}
}
