package main

import (
	"casbin_demo/middleware"
	"casbin_demo/pkg/db"
	"casbin_demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	userService := service.NewUserService(db.GetDB())
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", userService.Login)
	}

	productService := service.NewProductService(db.GetDB())
	productGroup := r.Group("/products")
	productGroup.Use(middleware.AuthMiddleware())
	{
		productGroup.POST("/add", productService.AddProductService)
		productGroup.DELETE("/del", productService.DeleteProductService)
		productGroup.POST("/upd", productService.UpdateProductService)
		productGroup.GET("/get", productService.GetProductService)
	}

	err = r.Run(":8811")
}
