package rocketmq

import (
	"auxiliary/rocketmq/types"
	"context"
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"log"
	"os"
	"time"
)

type SimpleCustomer struct{}

func NewCustomer() Consumer {
	return &SimpleCustomer{}
}

func (sc *SimpleCustomer) GetConsumer(ctx context.Context, c types.Conf) (*MqSubscribe, error) {
	var (
		d   = &MqSubscribe{}
		err error
	)

	err = os.Setenv("mq.consoleAppender.enabled", "true")
	if err != nil {
		return nil, err
	}
	rmqClient.ResetLogger()

	consumer, err := rmqClient.NewSimpleConsumer(&rmqClient.Config{
		Endpoint:      c.Endpoint,
		ConsumerGroup: c.TopicConf.Default.ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    c.AccessKey,
			AccessSecret: c.SecretKey,
		},
	},
		rmqClient.WithAwaitDuration(time.Second*5),
		rmqClient.WithSubscriptionExpressions(map[string]*rmqClient.FilterExpression{
			c.TopicConf.Default.Topic: rmqClient.NewFilterExpression("*"),
		}),
	)
	if err != nil {
		log.Fatalf("rocketmq consumer GetDefaultConsumer err: %v", err)
		return nil, err
	}

	d.reader = consumer
	d.conf = c
	d.ctx = ctx
	return d, nil
}

func (sc *SimpleCustomer) Subscribe(d *MqSubscribe) error {
	err := d.reader.Start()
	if err != nil {
		log.Fatalf("RocketMqSubscribe Start error:%v", err)
		return err
	}
	err = d.reader.Subscribe(d.conf.TopicConf.Default.Topic, rmqClient.NewFilterExpression("*"))
	if err != nil {
		log.Fatalf("RocketMqSubscribe Subscribe err: %v", err)
		return err
	}
	return nil
}

func (sc *SimpleCustomer) Stop(d *MqSubscribe) {
	func(simpleConsumer rmqClient.SimpleConsumer) {
		err := simpleConsumer.GracefulStop()
		if err != nil {
			log.Fatalf("RocketMqSubscribe GracefulStop() error:%v", err)
			return
		}
	}(d.reader)
}

func (sc *SimpleCustomer) Ack(d *MqSubscribe, mv *rmqClient.MessageView) error {
	err := d.reader.Ack(d.ctx, mv)
	if err != nil {
		log.Fatalf("RocketMqSubscribe ack error:%v", err)
		return err
	}
	return nil
}

func (sc *SimpleCustomer) Receive(d *MqSubscribe, maxMessageNum int32, invisibleDuration time.Duration) ([]*rmqClient.MessageView, error) {
	mvs, err := d.reader.Receive(d.ctx, maxMessageNum, invisibleDuration)
	if err != nil {
		log.Fatalf("RocketMqSubscribe Receive error:%v", err)
		return nil, err
	}
	return mvs, nil
}
