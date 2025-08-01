package main

import (
	"code_test/big_data/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func main() {

	db := model.InitDB()
	go transactionA(db)
	go transactionB(db)

	time.Sleep(15 * time.Second)
}

func transactionA(db *gorm.DB) {
	tx := db

	var ordersA []*model.Order
	var ordersB []*model.Order
	err := tx.Where("user_id = ? AND created_at >= ?", 23, "2025-03-17 00:00:00").Find(&ordersA).Error
	if err != nil {
		fmt.Printf("query orders error %s\n", err)
		return
	}
	fmt.Printf("query ordersA len =  %v\n", len(ordersA))
	time.Sleep(5 * time.Second) // 等待 TxB 插入

	err = tx.Where("user_id = ? AND created_at >= ?", 23, "2025-03-17 00:00:00").Find(&ordersB).Error
	if err != nil {
		fmt.Printf("query orders error %s\n", err)
		return
	}
	fmt.Printf("query ordersB len =  %v\n", len(ordersB))
}
func transactionB(db *gorm.DB) {
	time.Sleep(1 * time.Second)
	tx := db.Debug().Begin()
	defer tx.Rollback()

	var order = &model.Order{
		UserID:            23,
		UserName:          "Jack Baker",
		ProductID:         329,
		ProductName:       "思考，快与慢",
		ProductSpecsID:    407,
		ProductSpecsValue: "第2版",
		Quantity:          4,
		TotalPrice:        1196,
		OrderStatus:       3,
		Version:           0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	create := tx.Create(&order)
	if create.Error != nil {
		fmt.Printf("create order error %s\n", create.Error)
		return
	}
	if create.RowsAffected == 0 {
		fmt.Printf("create order RowsAffected ：0 \n")
		return
	}

	fmt.Printf("txb create success\n")
	tx.Commit()
}
