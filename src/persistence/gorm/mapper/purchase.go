package mapper

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/purchase"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models/quota"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/persistence/gorm/entities"
)

func ToPurchaseSinglePayment(entity *entities.PurchaseSinglePaymentEntity) *purchase.PurchaseSinglePayment {
	return &purchase.PurchaseSinglePayment{
		Purchase:      *toPurchase(&entity.PurchaseEntity),
		StoreDiscount: entity.StoreDiscount,
	}
}

func ToPurchaseSinglePaymentEntity(model *purchase.PurchaseSinglePayment) *entities.PurchaseSinglePaymentEntity {
	return &entities.PurchaseSinglePaymentEntity{
		PurchaseEntity: *ToPurchaseEntity(&model.Purchase),
		StoreDiscount:  model.StoreDiscount,
	}
}

func ToPurchaseMonthlyPayments(entity *entities.PurchaseMonthlyPaymentsEntity) *purchase.PurchaseMonthlyPayment {
	var quotas []quota.Quota
	for _, src := range entity.Quotas {
		quotas = append(quotas, *ToQuota(&src))
	}

	return &purchase.PurchaseMonthlyPayment{
		Purchase:       *toPurchase(&entity.PurchaseEntity),
		NumberOfQuotas: entity.NumberOfQuotas,
		Interest:       entity.Interest,
		Quota:          quotas,
	}
}

func ToPurchaseMonthlyPaymentsEntity(model *purchase.PurchaseMonthlyPayment) *entities.PurchaseMonthlyPaymentsEntity {
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

func ToPurchaseEntity(model *purchase.Purchase) *entities.PurchaseEntity {
	return &entities.PurchaseEntity{
		PaymentVoucher: model.PaymentVoucher,
		Store:          model.Store,
		CuitStore:      model.CuitStore,
		Amount:         model.Amount,
		FinalAmount:    model.FinalAmount,
	}
}

func toPurchase(entity *entities.PurchaseEntity) *purchase.Purchase {
	return &purchase.Purchase{
		PaymentVoucher: entity.PaymentVoucher,
		Store:          entity.Store,
		CuitStore:      entity.CuitStore,
		Amount:         entity.Amount,
		FinalAmount:    entity.FinalAmount,
	}
}

func ConvertPurchaseMonthlyPaymentsList(paymentEntityList *[]entities.PurchaseMonthlyPaymentsEntity) *[]purchase.PurchaseMonthlyPayment {
	var payments []purchase.PurchaseMonthlyPayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseMonthlyPayments(&v))
	}
	return &payments
}

func ConvertPurchaseSinglePaymentList(paymentEntityList *[]entities.PurchaseSinglePaymentEntity) *[]purchase.PurchaseSinglePayment {
	var payments []purchase.PurchaseSinglePayment
	for _, v := range *paymentEntityList {
		payments = append(payments, *ToPurchaseSinglePayment(&v))
	}
	return &payments
}
