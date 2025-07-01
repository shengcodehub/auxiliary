package mqtt

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/shengcodehub/auxiliary/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func GenerateSignature(accessKeySecret, clientID string) string {
	mac := hmac.New(sha1.New, []byte(accessKeySecret))
	mac.Write([]byte(clientID))
	signature := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func createClient(mqttConf *Conf) (MQTT.Client, error) {
	clientID := fmt.Sprintf("%s@@@%s", mqttConf.Group, "gg_push_"+os.Getenv("CLIENT_ID")+"_"+utils.RandomString(3)) // 自定义客户端ID
	// 生成签名
	signature := GenerateSignature(mqttConf.AccessKeySecret, clientID)
	fmt.Println("clientID:", clientID, "signature:", signature)
	// 创建 MQTT 客户端选项
	opts := MQTT.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:1883", mqttConf.Endpoint)).
		SetClientID(clientID).
		SetUsername(fmt.Sprintf("Signature|%s|%s", mqttConf.AccessKeyId, mqttConf.InstanceId)).
		SetPassword(signature).
		SetProtocolVersion(4)
	client := MQTT.NewClient(opts)
	// 连接到 MQTT 代理
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return nil, token.Error()
	}
	return client, nil
}

func SendMqttMessage(topic string, content string) (err error) {
	// 发布消息
	token := GetClient().Publish(topic, 2, false, content)
	res := token.Wait()
	logx.Infof("SendMqttMessage token.wait:%v, token.error:%v,  topic:%s, content:%s", res, token.Error(), topic, content)
	return token.Error()
}
