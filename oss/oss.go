package oss

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logc"
	"os"
)

type Conf struct {
	NameSrvAddr string
	AccessKey   string
	SecretKey   string
}
type BucketConf struct {
	Name string
	Url  string
}

func GetOss(c Conf, b BucketConf, cxt context.Context, fileName string, origin string, picData []byte) (string, error) {
	imageURL := origin + "/" + fileName
	client, err := oss.New(c.NameSrvAddr, c.AccessKey, c.SecretKey)
	if err != nil {
		logc.Errorf(cxt, fmt.Sprintf("GetOss New fail:%s", err.Error()))
		return "", err
	}
	bucket, err := client.Bucket(b.Name)
	if err != nil {
		logc.Errorf(cxt, fmt.Sprintf("GetOss Bucket fail:%s", err.Error()))
		return "", err
	}
	isExist, err := bucket.IsObjectExist(origin + "/" + fileName)
	if err != nil {
		logc.Errorf(cxt, fmt.Sprintf("GetOss IsObjectExist fail:%s", err.Error()))
		return "", err
	}
	if isExist {
		return imageURL, nil
	}
	path := "./" + fileName
	err = os.WriteFile(path, picData, 0644)
	if err != nil {
		logc.Errorf(cxt, fmt.Sprintf("GetOss WriteFile fail:%s", err.Error()))
		return "", err
	}
	err = bucket.PutObjectFromFile(origin+"/"+fileName, path)
	if err != nil {
		logc.Errorf(cxt, fmt.Sprintf("GetOss PutObjectFromFile fail:%s", err.Error()))
		return "", err
	}
	err = os.Remove(path)
	if err != nil {
		logc.Errorf(cxt, fmt.Sprintf("GetOss Remove fail:%s", err.Error()))
	}
	return imageURL, nil
}
