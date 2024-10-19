package customer

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

type Service interface {
	GetCustomerById(id int) (models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
}

type service struct {
	repo IRepository
}

func NewCustomerService(repo IRepository) Service {
	return &service{
		repo: repo,
	}
}

func (srvce *service) GetCustomerById(id int) (models.Customer, error) {
	customer, err := srvce.repo.GetCustomerById(id)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (srvce *service) GetAllCustomers() ([]models.Customer, error) {
	customers, err := srvce.repo.GetAllCustomers()
	if err != nil {
		return nil, err
	}
	return customers, nil
}
