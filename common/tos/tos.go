package tos

import (
	"context"
	"fmt"
	"github.com/shengcodehub/auxiliary/utils"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"
)

type Conf struct {
	Endpoint        string
	Region          string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          BucketConf
}

type BucketConf struct {
	Default string
}

var tosClient *tos.ClientV2

func SetUp(c Conf) {
	client, err := tos.NewClientV2(c.Endpoint, tos.WithRegion(c.Region), tos.WithCredentials(tos.NewStaticCredentials(c.AccessKeyId, c.AccessKeySecret)))
	if err != nil {
		logx.Errorf("tos NewClientV2 error:%v", err)
		return
	}
	tosClient = client
}

func GetClient() *tos.ClientV2 {
	return tosClient
}

func FileFormUpload(cxt context.Context, c Conf, r *http.Request, path string) (string, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		logx.Errorf("tos FileFormUpload error:%v", err)
		return "", err
	}
	arr := strings.Split(header.Filename, ".")
	var fileName string
	if len(arr) == 2 {
		fileName = utils.MD5(utils.GenerateUniqueIdentifier()) + "." + arr[1]
	} else {
		fileName = header.Filename
	}
	key := path + "/" + fileName
	upload := &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: c.Bucket.Default,
			Key:    key,
		},
		Content: file,
	}
	_, err = GetClient().PutObjectV2(cxt, upload)
	if err != nil {
		logx.Errorf("tos FileFormUpload PutObjectV2 error:%v", err)
		return "", err
	}
	return GetDefaultUrl(c) + key, nil
}

func GetDefaultUrl(c Conf) string {
	url := fmt.Sprintf("https://%s.%s/", c.Bucket.Default, c.Endpoint)
	return url
}
