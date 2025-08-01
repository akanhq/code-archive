package repositories

import (
	"code_test/message_queue/mongodb_database/blog-system/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByPostID(postID primitive.ObjectID, page, limit int) ([]*models.Comment, error)
	BatchCreate(comments []*models.Comment) error // 新增方法
}

type MongoCommentRepository struct {
	Collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *MongoCommentRepository {
	return &MongoCommentRepository{
		Collection: db.Collection("comments"),
	}
}

func (r *MongoCommentRepository) Create(comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.TODO(), comment)
	return err
}

func (r *MongoCommentRepository) FindByPostID(postID primitive.ObjectID, page, limit int) ([]*models.Comment, error) {

	skip := (page - 1) * limit
	filter := bson.M{"post_id": postID}
	cursor, err := r.Collection.Find(context.TODO(), filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "create_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	var comments []*models.Comment
	err = cursor.All(context.TODO(), &comments)
	return comments, err
}
func (r *MongoCommentRepository) BatchCreate(comments []*models.Comment) error {
	docs := make([]interface{}, len(comments))
	for i, comment := range comments {
		comment.CreatedAt = time.Now()
		docs[i] = comment
	}
	_, err := r.Collection.InsertMany(context.TODO(), docs)
	return err
}
