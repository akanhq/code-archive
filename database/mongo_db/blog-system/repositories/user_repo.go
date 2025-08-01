package repositories

import (
	"code_test/message_queue/mongodb_database/blog-system/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	BatchCreate(users []*models.User) error
}

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		Collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.TODO(), user)
	return err
}

func (r *MongoUserRepository) BatchCreate(users []*models.User) error {
	docs := make([]interface{}, len(users))
	for i, user := range users {
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		docs[i] = user
	}
	_, err := r.Collection.InsertMany(context.TODO(), docs)
	return err
}
func (r *MongoUserRepository) FindByUsername(username string) (*models.User, error) {
	filter := bson.M{"username": username}

	var user models.User
	err := r.Collection.FindOne(context.TODO(), filter).Decode(&user)
	return &user, err
}
