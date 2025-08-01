package model

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Order 定义订单表结构
type Order struct {
	ID                uint64    `gorm:"primaryKey"` // 订单ID
	UserID            uint64    // 用户ID
	UserName          string    // 用户姓名
	ProductID         uint64    // 商品ID
	ProductName       string    // 商品名称
	ProductSpecsID    uint64    // 商品规格ID
	ProductSpecsValue string    // 商品规格值
	Quantity          float64   // 商品数量
	TotalPrice        float64   // 订单总金额
	OrderStatus       int       // 订单状态 (0:待支付, 1:已支付, 2:已发货, 3:已完成, 4:已取消)
	Version           int       //乐观锁
	CreatedAt         time.Time // 创建时间
	UpdatedAt         time.Time // 更新时间
}

// 订单状态常量
const (
	StatusPending   = 0
	StatusPaid      = 1
	StatusShipped   = 2
	StatusCompleted = 3
	StatusCancelled = 4
)

func InitDB() *gorm.DB {
	dsn := "root:Xh123@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
	return db
}
