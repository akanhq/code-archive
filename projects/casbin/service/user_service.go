package service

import (
	"casbin_demo/models"
	"casbin_demo/pkg/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) Login(c *gin.Context) {
	var err error
	var req models.User

	if err = c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err = u.db.First(&user, req.Id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Id != req.Id || user.Username != req.Username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的用户信息"})
		return
	}

	tk := token.GenerateToken(&user)
	c.JSON(http.StatusOK, gin.H{"data": user, "token": tk})
}
