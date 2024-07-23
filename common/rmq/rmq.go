package rmq

import (
	"context"
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type Conf struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	TopicConf TopicConfMap
}

type TopicConfMap struct {
	Default TopicConf // 区分不同的topic
}

type TopicConf struct {
	Topic         string
	ConsumerGroup string
}

type MqBody struct {
	Event string `json:"event"` // 事件 用于区分走不同的逻辑
	Data  string `json:"data"`  // 数据
}

var producerClient rmqClient.Producer

func SetUp(c Conf) {
	err := os.Setenv("mq.consoleAppender.enabled", "true")
	if err != nil {
		logx.Errorf("rocketmq set env err: %v", err)
	}
	rmqClient.ResetLogger()
	producer, err := rmqClient.NewProducer(&rmqClient.Config{
		Endpoint: c.Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    c.AccessKey,
			AccessSecret: c.SecretKey,
		},
	})
	if err != nil {
		logx.Errorf("rocketmq producer err: %v", err)
	}

	err = producer.Start()
	if err != nil {
		logx.Errorf("rocketmq producer start err: %v", err)
	}
	producerClient = producer
}

func GetProduct() rmqClient.Producer {
	return producerClient
}

func GetDefaultProduct(ctx context.Context, c Conf, body *MqBody, keys []string, tag string) ([]*rmqClient.SendReceipt, error) {
	toString, err := jsoniter.Marshal(body)
	if err != nil {
		logx.Errorf("rocketmq producer jsoniter MarshalToString  err: %v", err)
		return nil, err
	}
	msg := &rmqClient.Message{
		Topic: c.TopicConf.Default.Topic,
		Body:  toString,
	}
	msg.SetKeys(keys...)
	msg.SetTag(tag)
	resp, err := GetProduct().Send(ctx, msg)
	if err != nil {
		logx.Errorf("rocketmq producer send err: %v", err)
		return nil, err
	}
	return resp, nil
}
