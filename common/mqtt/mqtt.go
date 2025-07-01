package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

type Conf struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	InstanceId      string
	Group           string
}

var MqttClient MQTT.Client

func SetUp(c *Conf) {
	client, err := createClient(c)
	if err != nil {
		logx.Errorf("mqtt client create err: %v", err)
	}
	MqttClient = client
}

func GetClient() MQTT.Client {
	return MqttClient
}
