package main

import (
	"code_test/logger"
	"log"

	"code_test/message_queue/mongodb_database/blog-system/config"
	"code_test/message_queue/mongodb_database/blog-system/handlers"
	"code_test/message_queue/mongodb_database/blog-system/middleware"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"code_test/message_queue/mongodb_database/blog-system/services"

	"github.com/gin-gonic/gin"
)

func main() {

	logger.InitLogger("debug")

	if err := config.InitDb(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	r := gin.Default()

	// 初始化服务和处理器
	userRepo := repositories.NewUserRepository(config.GetBlogDB())
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	postRepo := repositories.NewPostRepository(config.GetBlogDB())
	postService := services.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	commentRepo := repositories.NewCommentRepository(config.GetBlogDB())
	commentService := services.NewCommentService(commentRepo)
	commentHandler := handlers.NewCommentHandler(commentService)

	// 路由
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/posts", postHandler.CreatePost)
		api.GET("/posts", postHandler.GetPosts)
		api.GET("/posts/:id", postHandler.GetPostByID)
		api.PUT("/posts/:id", postHandler.UpdatePost)
		api.DELETE("/posts/:id", postHandler.DeletePost)
		api.GET("/posts/aggregate", postHandler.AggregatePosts)

		api.POST("/posts/:id/comments", commentHandler.CreateComment)
		api.GET("/posts/:id/comments", commentHandler.GetComments)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
