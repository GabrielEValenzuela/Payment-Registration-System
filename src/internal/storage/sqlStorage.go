package storage

import (
	"database/sql"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
)

type sqlStorage struct {
	DB *sql.DB
}

var QUERY_MAP = map[string]string{
	"GET_CUSTOMER":      "SELECT * FROM CUSTOMERS WHERE ID = ?",
	"GET_ALL_CUSTOMERS": "SELECT * FROM CUSTOMERS",
}

func NewSqlStorage(db *sql.DB) IStorage {
	return &sqlStorage{
		DB: db,
	}
}

func (ssql *sqlStorage) GetCustomerById(id int) (models.Customer, error) {
	var customer models.Customer
	err := ssql.DB.QueryRow(QUERY_MAP["GET_CUSTOMER"], id).Scan(&customer.CompleteName, &customer.Dni, &customer.Cuit, &customer.Address, &customer.Telephone, &customer.EntryDate)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (ssql *sqlStorage) GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	rows, err := ssql.DB.Query(QUERY_MAP["GET_CUSTOMER"])
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var customer models.Customer
		err = rows.Scan(&customer.CompleteName, &customer.Dni, &customer.Cuit, &customer.Address, &customer.Telephone, &customer.EntryDate)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}
