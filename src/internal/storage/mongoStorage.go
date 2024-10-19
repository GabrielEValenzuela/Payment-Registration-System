package storage

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoStorage struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStorage(client *mongo.Client, db *mongo.Database) IStorage {
	return &mongoStorage{
		client: client,
		db:     db,
	}
}

func (snsql *mongoStorage) GetCustomerById(id int) (models.Customer, error) {
	// Placeholder implementation
	return models.Customer{}, nil
}

func (snsql *mongoStorage) GetAllCustomers() ([]models.Customer, error) {
	// Placeholder implementation
	return nil, nil
}
