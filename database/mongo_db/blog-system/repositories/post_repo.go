package repositories

import (
	"code_test/logger"
	"code_test/message_queue/mongodb_database/blog-system/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id primitive.ObjectID) (*models.Post, error)
	FindAll(filter interface{}, page, limit int, sortField string, sortOrder int) ([]*models.Post, error)
	Update(id primitive.ObjectID, post *models.Post) error
	Delete(id primitive.ObjectID) error
	Aggregate(pipeline mongo.Pipeline) ([]map[string]interface{}, error)
}

type MongoPostRepository struct {
	Collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *MongoPostRepository {
	collection := db.Collection("posts")

	//创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "author_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().SetName("author_created_idx"),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		logger.Errorf("创建索引失败 %v", err)
	}

	return &MongoPostRepository{Collection: collection}
}

func (r *MongoPostRepository) Create(post *models.Post) error {
	_, err := r.Collection.InsertOne(context.TODO(), post)
	return err
}

func (r *MongoPostRepository) FindByID(id primitive.ObjectID) (*models.Post, error) {
	filter := bson.M{"_id": id}
	post := &models.Post{}
	err := r.Collection.FindOne(context.TODO(), filter).Decode(post)
	return post, err
}

func (r *MongoPostRepository) FindAll(filter interface{}, page, limit int, sortField string, sortOrder int) ([]*models.Post, error) {
	skip := (page - 1) * limit
	sort := bson.D{
		{Key: sortField},
		{Value: sortOrder},
	}
	cursor, err := r.Collection.Find(context.TODO(), filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var posts []*models.Post
	err = cursor.All(context.TODO(), &posts)
	return posts, err
}

func (r *MongoPostRepository) Update(id primitive.ObjectID, post *models.Post) error {
	filter := bson.M{"_id": id}
	post.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":     post.Title,
			"content":   post.Content,
			"update_at": post.UpdatedAt,
		},
	}
	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoPostRepository) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.Collection.DeleteOne(context.TODO(), filter)
	return err
}

func (r *MongoPostRepository) Aggregate(pipeline mongo.Pipeline) ([]map[string]interface{}, error) {
	cursor, err := r.Collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []map[string]interface{}
	err = cursor.All(context.TODO(), &results)
	return results, err
}
