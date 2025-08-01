package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (p *ProductService) AddProductService(c *gin.Context) {
	c.JSON(0, map[string]interface{}{
		"msg": "add product service",
	})
}

func (p *ProductService) DeleteProductService(c *gin.Context) {
	c.JSON(0, map[string]interface{}{
		"msg": "del product service",
	})
}

func (p *ProductService) UpdateProductService(c *gin.Context) {
	c.JSON(0, map[string]interface{}{
		"msg": "upd product service",
	})
}

func (p *ProductService) GetProductService(c *gin.Context) {
	c.JSON(0, map[string]interface{}{
		"msg": "get product service",
	})
}
