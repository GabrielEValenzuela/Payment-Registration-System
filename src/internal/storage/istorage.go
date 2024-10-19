package storage

import "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"

type IStorage interface {
	/* Customer */
	GetCustomerById(id int) (models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
	/* Bank */
}
