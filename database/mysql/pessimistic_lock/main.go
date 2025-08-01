package main

import (
	"code_test/big_data/model"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"gorm.io/gorm/clause"
)

func main() {
	db := model.InitDB()

	go transaction1(db)
	go transaction2(db)

	time.Sleep(10 * time.Second)
}

func transaction1(db *gorm.DB) {
	tx := db.Begin()
	if tx.Error != nil {
		fmt.Println("Tx1 begin error:", tx.Error)
		return
	}
	defer tx.Rollback()

	var order model.Order
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", 1).First(&order).Error
	if err != nil {
		fmt.Println("Tx1 query error:", err)
		return
	}
	nowTime := time.Now()
	fmt.Println("Tx1: Read total_price = ", order.TotalPrice, nowTime.Format("2006-01-02 15:04:05"))

	time.Sleep(5 * time.Second)

	order.TotalPrice += 50
	err = tx.Save(&order).Error
	if err != nil {
		fmt.Println("Tx1 save error:", err)
		return
	}
	tx.Commit()
	fmt.Println("Tx1: upd total_price = ", order.TotalPrice, nowTime.Format("2006-01-02 15:04:05"))
	fmt.Println("Tx1: Committed")
}

func transaction2(db *gorm.DB) {
	time.Sleep(500 * time.Millisecond)
	tx := db.Begin()
	if tx.Error != nil {
		fmt.Println("Tx2 begin error:", tx.Error)
		return
	}
	defer tx.Rollback()

	var order model.Order
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", 1).First(&order).Error
	if err != nil {
		fmt.Println("Tx2 query error:", err)
		return
	}
	nowTime := time.Now()
	fmt.Println("Tx2: Read total_price = ", order.TotalPrice, nowTime.Format("2006-01-02 15:04:05"))

	order.TotalPrice = 300
	err = tx.Save(&order).Error
	if err != nil {
		log.Println("Tx2 update error:", err)
		return
	}

	nowTime = time.Now()
	fmt.Println("Tx2: upd after total_price = ", order.TotalPrice, nowTime.Format("2006-01-02 15:04:05"))

	tx.Commit()
	fmt.Println("Tx2: Committed")
}
