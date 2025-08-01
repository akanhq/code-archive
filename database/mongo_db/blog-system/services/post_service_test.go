package services

import (
	"code_test/message_queue/mongodb_database/blog-system/config"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"code_test/message_queue/mongodb_database/blog-system/models"
	"go.mongodb.org/mongo-driver/bson"
)

const postCount = 10000000
const postSize = postCount / workerCount

func TestBatchCreatePosts(t *testing.T) {
	start := time.Now()

	if err := config.InitDb(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	userRepo := repositories.NewUserRepository(config.GetBlogDB())

	// 获取用户列表
	cursor, err := userRepo.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Printf("获取用户列表失败,err:%v\n", err)
		return
	}
	var users []*models.User
	err = cursor.All(context.TODO(), &users)

	var wg sync.WaitGroup
	wg.Add(postCount)

	var batchSize = 1000
	for workID := 0; workID < workerCount; workID++ {
		go func() {
			defer wg.Done()
			for i := 0; i < batchSize; i++ {

			}
		}()
	}
	wg.Wait()

	fmt.Printf("Inserted %d posts in %v\n", postCount, time.Since(start))
}
