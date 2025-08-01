package db

import (
	"casbin_demo/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
)

var DB *gorm.DB

const (
	SUPER_ID        = 9999
	SUPER_USER_NAME = "SuperAdmin"
)

func init() {
	// 主库 DSN
	masterDSN := "root:Xh123@tcp(14.103.142.181:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
	// 从库 DSN
	slaveDSN := "root:Xh123@tcp(14.103.142.181:3307)/database?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接主库
	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to master database: %v", err)
	}

	// 配置读写分离
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(masterDSN)}, // 主库用于写
		Replicas: []gorm.Dialector{mysql.Open(slaveDSN)},  // 从库用于读
		Policy:   dbresolver.RandomPolicy{},               // 随机选择从库（支持多个从库）
	}))
	if err != nil {
		log.Fatalf("Failed to configure db resolver: %v", err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.Permission{})
	db.AutoMigrate(&models.RolePermission{})

	//用户
	var superAdmin models.User
	if err = db.First(&superAdmin, SUPER_ID).Error; err != nil {
		superAdmin = models.User{
			CommonField: models.CommonField{Id: SUPER_ID},
			Username:    SUPER_USER_NAME,
		}
		db.Create(&superAdmin)
	}

	//角色
	var superRole models.Role
	if err = db.First(&superRole, SUPER_ID).Error; err != nil {
		superRole = models.Role{
			CommonField: models.CommonField{Id: SUPER_ID},
			Name:        "超级管理员权限",
		}
		db.Create(&superRole)
	}

	//权限
	var superPermission models.Permission
	if err = db.First(&superPermission, SUPER_ID).Error; err != nil {
		superPermission = models.Permission{
			CommonField: models.CommonField{Id: SUPER_ID},
			Name:        "*:*",
		}
		db.Create(&superPermission)
	}

	//绑定
	var rolePermission models.RolePermission
	if err = db.First(&rolePermission, SUPER_ID).Error; err != nil {
		rolePermission = models.RolePermission{
			CommonField:  models.CommonField{Id: SUPER_ID},
			RoleId:       SUPER_ID,
			PermissionId: SUPER_ID,
		}
		db.Create(&rolePermission)
	}

	//绑定超级用户到超级管理员角色
	superAdmin.RoleId = superRole.Id
	db.Save(&superAdmin)
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
