package main

import (
	"code_test/big_data/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func main() {

	db := model.InitDB()

	go optimisticUpdate(db, "Tx1", 50)
	go optimisticUpdate(db, "Tx2", 100)

	time.Sleep(10 * time.Second)
}

func optimisticUpdate(db *gorm.DB, name string, increment float64) {
	tx := db.Debug().Begin()
	if tx.Error != nil {
		fmt.Printf("%s begin error:%s\n", name, tx.Error)
	}
	defer tx.Rollback()

	// query data
	var order model.Order
	err := tx.Where("id = ?", 1).First(&order).Error
	if err != nil {
		fmt.Printf("%s order not exists: %s\n", name, err)
	}
	fmt.Printf("%s Read Version = %d && total_price = %.2f\n", name, order.Version, order.TotalPrice)

	// update data
	result := tx.Model(&model.Order{}).Where("id = ? AND version = ?", order.ID, order.Version).
		Updates(map[string]interface{}{
			"total_price": increment,
			"version":     order.Version + 1,
		})

	if result.Error != nil {
		fmt.Printf("%s update order fail:%s\n", name, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		fmt.Printf("%s Update failed, version changed by another transaction\n", name)
	} else {
		fmt.Printf("%s update order success\n", name)
	}

	tx.Commit()

}
