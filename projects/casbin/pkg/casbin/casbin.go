package casbin

import (
	"casbin_demo/pkg/db"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func InitCasbin() error {
	// 使用 GORM 适配器
	adapter, err := gormadapter.NewAdapterByDB(db.GetDB())
	if err != nil {
		return err
	}

	// 使用基础的 RBAC 模型
	// model.conf 文件你可以放在项目根目录或 config 目录
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	Enforcer = enforcer
	return nil
}
