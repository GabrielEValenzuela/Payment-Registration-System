package storage

import (
	"database/sql"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

type sqlStorage struct {
	DB *sql.DB
}

var QUERY_MAP = map[string]string{
	"GET_CUSTOMER": "SELECT * FROM CUSTOMERS WHERE ID = ?",
}

func NewSqlStorage(db *sql.DB) IStorage {
	return &sqlStorage{
		DB: db,
	}
}

func (s *sqlStorage) GetCustomerById(id int) (models.Customer, error) {
	var customer models.Customer
	err := s.DB.QueryRow(QUERY_MAP["GET_CUSTOMER"], id).Scan(&customer.CompleteName, &customer.Dni, &customer.Cuit, &customer.Address, &customer.Telephone, &customer.EntryDate)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}
