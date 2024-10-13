package customer

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
)

type IRepository interface {
	GetCustomerById(id int) (models.Customer, error)
}

type repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &repository{
		storage: storage,
	}
}

func (r *repository) GetCustomerById(id int) (models.Customer, error) {
	return r.storage.GetCustomerById(id)
}
