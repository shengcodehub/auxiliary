package orm

import (
	"github.com/shengcodehub/auxiliary/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	_db *gorm.DB
)

func Setup(dataSource string) {
	var (
		err error
		url = dataSource
	)

	_db, err = gorm.Open(mysql.Open(url))
	if err != nil {
		log.Fatalf("open gorm failed %s", err.Error())
	}

	_, err = _db.DB()
	if err != nil {
		log.Fatalf("get gorm db failed %s", err.Error())
	}

	if utils.IsDev() || utils.IsTest() {
		_db = _db.Debug()
	}
}

func GetDB() *gorm.DB {
	return _db
}
