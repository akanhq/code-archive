package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitDb() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10秒超时
	defer cancel()

	clientOptions := options.Client().
		ApplyURI("mongodb://Xh123:Dev123@localhost:27017/?authSource=admin").
		SetMaxPoolSize(20). // 连接池大小
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Second) // 空闲连接超时

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	log.Println("Connected to MongoDB!")
	return nil
}

func GetBlogDB() *mongo.Database {
	return client.Database("blog_db")
}

func CloseDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}
