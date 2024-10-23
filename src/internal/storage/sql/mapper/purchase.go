package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/sql/entities"
)

func ToPurchaseSinglePayment(entity *entities.PurchaseSinglePaymentEntity) *models.PurchaseSinglePayment {
	return &models.PurchaseSinglePayment{
		Purchase:      *toPurchase(&entity.PurchaseEntity),
		StoreDiscount: entity.StoreDiscount,
	}
}

func ToPurchaseSinglePaymentEntity(model *models.PurchaseSinglePayment) *entities.PurchaseSinglePaymentEntity {
	return &entities.PurchaseSinglePaymentEntity{
		PurchaseEntity: *ToPurchaseEntity(&model.Purchase),
		StoreDiscount:  model.StoreDiscount,
	}
}

func ToPurchaseMonthlyPayments(entity *entities.PurchaseMonthlyPaymentsEntity) *models.PurchaseMonthlyPayment {
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

func ToPurchaseMonthlyPaymentsEntity(model *models.PurchaseMonthlyPayment) *entities.PurchaseMonthlyPaymentsEntity {
	var quotas []entities.QuotaEntity
	for _, src := range model.Quota {
		quotas = append(quotas, *ToQuotaEntity(&src))
	}
	return &entities.PurchaseMonthlyPaymentsEntity{
		PurchaseEntity: *ToPurchaseEntity(&model.Purchase),
		Interest:       model.Interest,
		NumberOfQuotas: model.NumberOfQuotas,
		Quotas:         quotas,
	}
}

func ToPurchaseEntity(model *models.Purchase) *entities.PurchaseEntity {
	return &entities.PurchaseEntity{
		PaymentVoucher: model.PaymentVoucher,
		Store:          model.Store,
		CuitStore:      model.CuitStore,
		Amount:         model.Amount,
		FinalAmount:    model.FinalAmount,
	}
}

func toPurchase(entity *entities.PurchaseEntity) *models.Purchase {
	return &models.Purchase{
		PaymentVoucher: entity.PaymentVoucher,
		Store:          entity.Store,
		CuitStore:      entity.CuitStore,
		Amount:         entity.Amount,
		FinalAmount:    entity.FinalAmount,
	}
}

func ConvertPurchaseMonthlyPaymentsList(paymentEntityList *[]entities.PurchaseMonthlyPaymentsEntity) *[]models.PurchaseMonthlyPayment {
	var payments []models.PurchaseMonthlyPayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseMonthlyPayments(&v))
	}
	return &payments
}

func ConvertPurchaseSinglePaymentList(paymentEntityList *[]entities.PurchaseSinglePaymentEntity) *[]models.PurchaseSinglePayment {
	var payments []models.PurchaseSinglePayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseSinglePayment(&v))
	}
	return &payments
}
