package services

import (
	"code_test/message_queue/mongodb_database/blog-system/models"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	Repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) *CommentService {
	return &CommentService{Repo: repo}
}

func (s *CommentService) CreateComment(comment *models.Comment) error {
	return s.Repo.Create(comment)
}

func (s *CommentService) GetCommentsByPostID(postID primitive.ObjectID, page, limit int) ([]*models.Comment, error) {
	return s.Repo.FindByPostID(postID, page, limit)
}
func (s *CommentService) BatchCreateComments(comments []*models.Comment) error {
	return s.Repo.BatchCreate(comments)
}
