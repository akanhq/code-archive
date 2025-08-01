package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"code_test/message_queue/mongodb_database/blog-system/config"
	"code_test/message_queue/mongodb_database/blog-system/models"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupCommentTest(t *testing.T) (*CommentService, []*models.User, []*models.Post, func()) {
	if err := config.InitDb(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	db := config.GetBlogDB()
	commentRepo := repositories.NewCommentRepository(db)
	commentService := NewCommentService(commentRepo)
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)

	// 插入测试用户
	users := make([]*models.User, 0, 100)
	for i := 0; i < 100; i++ {
		user := &models.User{
			ID:        primitive.NewObjectID(),
			Username:  gofakeit.Username(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, true, false, 12),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)
		users = append(users, user)
	}

	// 插入测试文章
	posts := make([]*models.Post, 0, 10)
	for i := 0; i < 10; i++ {
		user := users[i%len(users)]
		post := &models.Post{
			ID:        primitive.NewObjectID(),
			Title:     gofakeit.Sentence(5),
			Content:   gofakeit.Paragraph(1, 3, 10, "."),
			AuthorID:  user.ID,
			Author:    user.Username,
			Views:     gofakeit.Number(0, 1000),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := postRepo.Create(post)
		assert.NoError(t, err)
		posts = append(posts, post)
	}

	// 清理函数
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		collections := []string{"users", "posts", "comments"}
		for _, col := range collections {
			_, err := db.Collection(col).DeleteMany(ctx, bson.M{})
			if err != nil {
				log.Printf("Failed to clean %s collection: %v", col, err)
			}
		}
		config.CloseDB()
	}

	return commentService, users, posts, cleanup
}

func TestBatchCreateComments(t *testing.T) {
	commentService, users, posts, cleanup := setupCommentTest(t)
	defer cleanup()

	const commentCount = 2_000_000
	const batchSize = 1000
	start := time.Now()

	// 使用并发批量插入评论
	var wg sync.WaitGroup
	jobs := make(chan []*models.Comment, 10)
	errChan := make(chan error, 10)

	// 启动 5 个工人 goroutine
	workerCount := 5
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range jobs {
				if err := commentService.BatchCreateComments(batch); err != nil {
					errChan <- err
					return
				}
			}
		}()
	}

	// 生成评论并分发
	batch := make([]*models.Comment, 0, batchSize)
	for i := 0; i < commentCount; i++ {
		user := users[i%len(users)]
		post := posts[i%len(posts)]
		comment := &models.Comment{
			ID:        primitive.NewObjectID(),
			PostID:    post.ID,
			AuthorID:  user.ID,
			Author:    user.Username,
			Content:   gofakeit.Sentence(10),
			CreatedAt: time.Now(),
		}
		batch = append(batch, comment)

		if len(batch) >= batchSize {
			jobs <- batch
			fmt.Printf("Generated %d comments\n", i+1)
			batch = make([]*models.Comment, 0, batchSize)
		}
	}

	if len(batch) > 0 {
		jobs <- batch
	}
	close(jobs)

	// 等待所有工人完成
	wg.Wait()
	close(errChan)

	// 检查错误
	for err := range errChan {
		assert.NoError(t, err, "Failed to insert comment batch")
	}

	// 验证插入数量
	count, err := commentService.Repo.(*repositories.MongoCommentRepository).Collection.CountDocuments(context.TODO(), bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(commentCount), count, "Expected %d comments, got %d", commentCount, count)

	// 验证随机评论数据
	var comment models.Comment
	err = commentService.Repo.(*repositories.MongoCommentRepository).Collection.FindOne(context.TODO(), bson.M{}).Decode(&comment)
	assert.NoError(t, err)
	assert.NotEmpty(t, comment.Content)
	assert.NotEmpty(t, comment.AuthorID)
	assert.NotEmpty(t, comment.PostID)

	fmt.Printf("Inserted %d comments in %v\n", commentCount, time.Since(start))
}
