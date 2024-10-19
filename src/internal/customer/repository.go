package customer

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
)

type IRepository interface {
	GetCustomerById(id int) (models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
}

type repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &repository{
		storage: storage,
	}
}

func (repo *repository) GetCustomerById(id int) (models.Customer, error) {
	return repo.storage.GetCustomerById(id)
}

func (repo *repository) GetAllCustomers() ([]models.Customer, error) {
	return repo.storage.GetAllCustomers()
}
