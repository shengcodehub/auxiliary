package sms

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type Conf struct {
	Endpoint     string
	AccessKey    string
	SecretKey    string
	SignName     string
	TemplateConf TemplateCodeConf
}

type TemplateCodeConf struct {
	ScheduleApply string
}

var aliyunClient *dysmsapi.Client

func SetUp(c Conf) {
	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(c.AccessKey),
		AccessKeySecret: tea.String(c.SecretKey),
		Endpoint:        tea.String(c.Endpoint),
	})
	if err != nil {
		logx.Errorf("dysmsapi NewClient error: %s", err.Error())
		return
	}
	aliyunClient = client
}

func SendSms(phones []string, signName string, templateCode string, templateParam string) error {
	sendSmsRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(strings.Join(phones, ",")),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(templateParam),
	}
	sms, err := aliyunClient.SendSms(sendSmsRequest)
	if err != nil {
		logx.Errorf("dysmsapi SendSms error: %s", err.Error())
		return err
	}
	response, _ := jsoniter.MarshalToString(sms.Body)
	logx.Infof("dysmsapi SendSms phones: %s , response: %s", strings.Join(phones, ","), response)
	if tea.StringValue(sms.Body.Code) != "ok" {
		return errors.New(tea.StringValue(sms.Body.Message))
	}
	return nil
}
