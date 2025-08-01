package services

import (
	"code_test/message_queue/mongodb_database/blog-system/config"
	"code_test/message_queue/mongodb_database/blog-system/models"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/bxcodec/faker/v3"
)

const userCount = 1000000
const workerCount = 10
const userSize = userCount / workerCount

func TestBatchCreateUsers(t *testing.T) {

	start := time.Now()

	if err := config.InitDb(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	userRepo := repositories.NewUserRepository(config.GetBlogDB())
	authService := NewAuthService(userRepo)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go batchCreateUsers(i, &wg, authService)
	}

	wg.Wait()
	fmt.Printf("Inserted %d comments in %v\n", userCount, time.Since(start))
}

func batchCreateUsers(workID int, wg *sync.WaitGroup, authService *AuthService) {
	defer wg.Done()

	var err error
	var users []*models.User
	var batchSize = 1000
	for i := 0; i < userSize; i++ {

		user := &models.User{
			Username: gofakeit.Username(),
			Password: gofakeit.Password(true, true, true, true, true, 12),
			Email:    gofakeit.Email(),
		}
		if i%2 == 0 {
			user.Username = faker.ChineseName()
		}
		users = append(users, user)

		if i%batchSize == 0 {
			err = authService.BatchCreate(users)
			users = make([]*models.User, 0)
			if err != nil {
				fmt.Printf("批量创建用户失败, %v\n", err)
				continue
			}
			fmt.Printf("工人 %d 完成数量 %d\n", workID, i)
		}
	}
}
