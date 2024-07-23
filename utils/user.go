package utils

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"strconv"
	"time"
)

func GetCity(iP string, geoip2File string) (*geoip2.City, error) {
	db, err := geoip2.Open(geoip2File)
	if err != nil {
		logx.Errorf("open geoip2 err:%s", err.Error())
		return nil, err
	}
	defer func(db *geoip2.Reader) {
		err := db.Close()
		if err != nil {
			logx.Errorf("close geoip2 err:%s", err.Error())
		}
	}(db)

	ip := net.ParseIP(iP)

	record, err := db.City(ip)
	if err != nil {
		logx.Errorf("get ip info err:%s", err.Error())
		return nil, err
	}
	return record, nil
}

// GetAgeByIDCard 根据身份证号码计算年龄
func GetAgeByIDCard(idCard string) (int, error) {
	if len(idCard) != 18 {
		return 0, fmt.Errorf("invalid ID card number, should be 18 digits")
	}

	// 提取出生日期部分
	birthYear := idCard[6:10]
	birthMonth := idCard[10:12]
	birthDay := idCard[12:14]
	birthDate := birthYear + "-" + birthMonth + "-" + birthDay

	// 解析出生日期
	birthTime, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, fmt.Errorf("failed to parse birth date: %w", err)
	}

	// 计算年龄
	age := int(time.Now().Sub(birthTime).Hours() / 24 / 365)

	return age, nil
}

// GetGenderByIDCard 性别 1 女 2 男
func GetGenderByIDCard(idCard string) (int64, error) {
	if len(idCard) != 18 {
		return 0, fmt.Errorf("invalid ID card number, should be 18 digits")
	}
	// 提取倒数第二位，即性别码
	genderCode := idCard[len(idCard)-2 : len(idCard)-1]
	genderNum, err := strconv.Atoi(genderCode)
	if err != nil {
		return 0, fmt.Errorf("无法解析性别码: %w", err)
	}

	// 根据性别码判断性别
	if genderNum%2 == 0 {
		return 1, nil
	}
	return 2, nil
}
