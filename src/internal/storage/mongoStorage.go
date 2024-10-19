package storage

import (
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
