package repositories

import (
	"code_test/message_queue/mongodb_database/contact-manager/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactRepository interface {
	Create(contact *models.Contact) error
	FindByID(id primitive.ObjectID) (*models.Contact, error)
	FindByName(name string) ([]*models.Contact, error)
	Update(id primitive.ObjectID, contact *models.Contact) error
	Delete(id primitive.ObjectID) error
}

type MongoContactRepository struct {
	Collection *mongo.Collection
}

func NewContactRepository(db *mongo.Database) *MongoContactRepository {
	return &MongoContactRepository{
		Collection: db.Collection("contacts"),
	}
}

func (r *MongoContactRepository) Create(contact *models.Contact) error {
	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.TODO(), contact)
	return err
}

func (r *MongoContactRepository) FindByID(id primitive.ObjectID) (*models.Contact, error) {
	filter := bson.M{"_id": id}
	var contact models.Contact
	err := r.Collection.FindOne(context.TODO(), filter).Decode(&contact)
	return &contact, err
}

func (r *MongoContactRepository) FindByName(name string) ([]*models.Contact, error) {
	filter := bson.M{"name": name}

	cursor, err := r.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var contacts []*models.Contact
	err = cursor.All(context.TODO(), &contacts)
	return contacts, err
}

func (r *MongoContactRepository) Update(id primitive.ObjectID, contact *models.Contact) error {
	contact.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":       contact.Name,
			"phone":      contact.Phone,
			"email":      contact.Email,
			"updated_at": contact.UpdatedAt,
		},
	}
	filter := bson.M{"_id": id}
	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}
func (r *MongoContactRepository) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.Collection.DeleteOne(context.TODO(), filter)
	return err
}
