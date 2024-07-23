package orm

import (
	"github.com/shengwenjin/auxiliary/common/env"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	_db *gorm.DB
)

func Setup() {
	setupDb()
}

func setupDb() {
	var (
		err error
		url = viper.GetString("Mysql.Rds.DataSource")
	)

	_db, err = gorm.Open(mysql.Open(url))
	if err != nil {
		logx.Errorf("open gorm failed %s", err.Error())
	}

	db, err := _db.DB()
	if err != nil {
		logx.Errorf("get gorm db failed %s", err.Error())
	}

	db.SetConnMaxLifetime(time.Duration(viper.GetInt("Mysql.Rds.ConnMaxLifetime")) * time.Second)
	db.SetMaxIdleConns(viper.GetInt("Mysql.Rds.MaxIdleConn"))
	db.SetMaxOpenConns(viper.GetInt("Mysql.Rds.MaxOpenConn"))

	if env.IsDev() {
		_db = _db.Debug()
	}
}

func GetDB() *gorm.DB {
	return _db
}
