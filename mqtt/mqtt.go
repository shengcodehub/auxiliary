package mqtt

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	onsmqtt20200420 "github.com/alibabacloud-go/onsmqtt-20200420/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type (
	AliyunMqtt interface {
		SendMqttMessage(topic string, content string) (err error)
		CreateDeviceToken(clientId string) (resp *onsmqtt20200420.RegisterDeviceCredentialResponse, err error)
		CreateMqttToken() (resp *onsmqtt20200420.ApplyTokenResponse, err error)
		UnRegisterDevice(clientId string) (err error)
	}

	aliyunMqttConfig struct {
		config *mqttConfig
	}

	mqttConfig struct {
		AccessKeyId     string
		AccessKeySecret string
		Endpoint        string
		InstanceId      string
		RegionId        string
	}
)

func NewAliyunMqtt(accessKeyId, accessKeySecret, endpoint, instanceId, regionId string) AliyunMqtt {
	return &aliyunMqttConfig{
		config: &mqttConfig{
			AccessKeyId:     accessKeyId,
			AccessKeySecret: accessKeySecret,
			Endpoint:        endpoint,
			InstanceId:      instanceId,
			RegionId:        regionId,
		},
	}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func createClient(accessKeyId *string, accessKeySecret *string, regionId *string) (_result *onsmqtt20200420.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/OnsMqtt
	config.Endpoint = tea.String(fmt.Sprintf("onsmqtt.%s.aliyuncs.com", *regionId))
	_result = &onsmqtt20200420.Client{}
	_result, _err = onsmqtt20200420.NewClient(config)
	return _result, _err
}

/**
* 使用STS鉴权方式初始化账号Client，推荐此方式。
* @param accessKeyId
* @param accessKeySecret
* @param securityToken
* @return Client
* @throws Exception
 */

func createClientWithSTS(accessKeyId *string, accessKeySecret *string, securityToken *string, regionId *string) (_result *onsmqtt20200420.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
		// 必填，您的 Security Token
		SecurityToken: securityToken,
		// 必填，表明使用 STS 方式
		Type: tea.String("sts"),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/OnsMqtt
	config.Endpoint = tea.String(fmt.Sprintf("onsmqtt.%s.aliyuncs.com", *regionId))
	_result = &onsmqtt20200420.Client{}
	_result, _err = onsmqtt20200420.NewClient(config)
	return _result, _err
}

func (m *aliyunMqttConfig) CreateDeviceToken(clientId string) (resp *onsmqtt20200420.RegisterDeviceCredentialResponse, err error) {
	client, err := createClient(tea.String(m.config.AccessKeyId), tea.String(m.config.AccessKeySecret), tea.String(m.config.RegionId))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	registerDeviceCredentialRequest := &onsmqtt20200420.RegisterDeviceCredentialRequest{
		InstanceId: tea.String(m.config.InstanceId),
		ClientId:   tea.String(clientId),
	}

	runtime := &util.RuntimeOptions{}
	resp, tryErr := func() (resp *onsmqtt20200420.RegisterDeviceCredentialResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, _err := client.RegisterDeviceCredentialWithOptions(registerDeviceCredentialRequest, runtime)
		if _err != nil {
			return nil, _err
		}

		return resp, nil
	}()

	if tryErr != nil {
		return nil, fmt.Errorf("failed to register device credential: %w", tryErr)
	}
	return resp, err
}

func (m *aliyunMqttConfig) SendMqttMessage(topic string, content string) (err error) {
	client, err := createClient(tea.String(m.config.AccessKeyId), tea.String(m.config.AccessKeySecret), tea.String(m.config.RegionId))
	if err != nil {
		return err
	}

	sendMessageRequest := &onsmqtt20200420.SendMessageRequest{
		Payload:    tea.String(content),
		InstanceId: tea.String(m.config.InstanceId),
		MqttTopic:  tea.String(topic),
	}
	runtime := &util.RuntimeOptions{}
	_, tryErr := func() (resp *onsmqtt20200420.SendMessageResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, err := client.SendMessageWithOptions(sendMessageRequest, runtime)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}()

	if tryErr != nil {
		return tryErr
	}

	return nil
}

func (m *aliyunMqttConfig) CreateMqttToken() (resp *onsmqtt20200420.ApplyTokenResponse, err error) {
	client, err := createClient(tea.String(m.config.AccessKeyId), tea.String(m.config.AccessKeySecret), tea.String(m.config.RegionId))

	if err != nil {
		return nil, err
	}

	applyTokenRequest := &onsmqtt20200420.ApplyTokenRequest{
		Actions:    tea.String("R"),
		ExpireTime: tea.Int64(1709222400000),
		InstanceId: tea.String("post-cn-7mz2e5hc90i"),
		Resources:  tea.String("system"),
	}
	runtime := &util.RuntimeOptions{}
	resp, tryErr := func() (resp *onsmqtt20200420.ApplyTokenResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, err = client.ApplyTokenWithOptions(applyTokenRequest, runtime)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}()

	if tryErr != nil {
		return nil, tryErr
	}
	return resp, err
}

func (m *aliyunMqttConfig) UnRegisterDevice(clientId string) (err error) {
	client, err := createClient(tea.String(m.config.AccessKeyId), tea.String(m.config.AccessKeySecret), tea.String(m.config.RegionId))
	if err != nil {
		return err
	}

	unRegisterDeviceCredentialRequest := &onsmqtt20200420.UnRegisterDeviceCredentialRequest{
		ClientId:   tea.String(clientId),
		InstanceId: tea.String(m.config.InstanceId),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, err = client.UnRegisterDeviceCredentialWithOptions(unRegisterDeviceCredentialRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		return tryErr
	}
	return err
}
