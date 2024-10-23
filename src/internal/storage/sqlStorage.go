package storage

import (
	"database/sql"
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
