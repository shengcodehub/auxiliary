package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var enforcer *casbin.Enforcer

func SetUpByDb(db *gorm.DB) {
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		logx.Errorf("casbin adapter db init error: %v", err)
		return
	}
	e, err := casbin.NewEnforcer("acl_model.conf", a)
	if err != nil {
		logx.Errorf("casbin db init error: %v", err)
		return
	}
	err = e.LoadPolicy()
	if err != nil {
		logx.Errorf("casbin db load policy error: %v", err)
		return
	}
	enforcer = e
}

func GetCasbin(sub, obj, act string) bool {
	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		logx.Errorf("casbin Enforce error: %v", err)
		return false
	}
	return ok
}

func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
