package customer

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

type Service interface {
	GetCustomerById(id int) (models.Customer, error)
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetCustomerById(id int) (models.Customer, error) {
	customer, err := s.repo.GetCustomerById(id)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}
