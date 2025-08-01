package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitDb() error {

	clientOptions := options.Client().ApplyURI("mongodb://Xh123:Dev123@localhost:27017/?authSource=admin")

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	log.Println("Connected to MongoDB!")
	return nil
}

func GetContactDB() *mongo.Database {
	return client.Database("contact_db")
}

func CloseDB() error {
	return client.Disconnect(context.TODO())
}
