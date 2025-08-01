package services

import (
	"code_test/message_queue/mongodb_database/blog-system/config"
	"context"

	"code_test/message_queue/mongodb_database/blog-system/models"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostService struct {
	Repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) *PostService {
	return &PostService{Repo: repo}
}

func (s *PostService) CreatePost(post *models.Post) error {
	return s.Repo.Create(post)
}

func (s *PostService) GetPostByID(id primitive.ObjectID) (*models.Post, error) {
	return s.Repo.FindByID(id)
}

func (s *PostService) GetPosts(filter map[string]interface{}, page, limit int, sortField, sortOrder string) ([]*models.Post, error) {
	sortOrderInt := 1
	if sortOrder == "desc" {
		sortOrderInt = -1
	}
	return s.Repo.FindAll(filter, page, limit, sortField, sortOrderInt)
}

func (s *PostService) UpdatePost(id primitive.ObjectID, post *models.Post) error {
	return s.Repo.Update(id, post)
}

func (s *PostService) DeletePost(id primitive.ObjectID) error {
	return s.Repo.Delete(id)
}

func (s *PostService) AggregatePosts(pipeline mongo.Pipeline) ([]map[string]interface{}, error) {
	return s.Repo.Aggregate(pipeline)
}

func (s *PostService) BatchCreatePosts(posts []*models.Post) error {
	session, err := config.GetBlogDB().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	ctx := context.TODO()
	for _, post := range posts {
		if err := s.Repo.Create(post); err != nil {
			session.AbortTransaction(ctx)
			return err
		}
	}

	return session.CommitTransaction(ctx)
}
