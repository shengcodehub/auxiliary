package rocketmq

import (
	"context"
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/shengcodehub/auxiliary/rocketmq/types"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("rocketmq consumer GetDefaultConsumer err: %v", err)
		return nil, err
	}

	d.Reader = consumer
	d.Conf = c
	d.Ctx = ctx
	return d, nil
}

func (sc *SimpleCustomer) Subscribe(d *MqSubscribe) error {
	err := d.Reader.Start()
	if err != nil {
		log.Errorf("RocketMqSubscribe Start error:%v", err)
		return err
	}
	err = d.Reader.Subscribe(d.Conf.TopicConf.Default.Topic, rmqClient.NewFilterExpression("*"))
	if err != nil {
		log.Errorf("RocketMqSubscribe Subscribe err: %v", err)
		return err
	}
	return nil
}

func (sc *SimpleCustomer) Stop(d *MqSubscribe) {
	func(simpleConsumer rmqClient.SimpleConsumer) {
		err := simpleConsumer.GracefulStop()
		if err != nil {
			log.Errorf("RocketMqSubscribe GracefulStop() error:%v", err)
			return
		}
	}(d.Reader)
}

func (sc *SimpleCustomer) Ack(d *MqSubscribe, mv *rmqClient.MessageView) error {
	err := d.Reader.Ack(d.Ctx, mv)
	if err != nil {
		log.Errorf("RocketMqSubscribe ack error:%v", err)
		return err
	}
	return nil
}

func (sc *SimpleCustomer) Receive(d *MqSubscribe, maxMessageNum int32, invisibleDuration time.Duration) ([]*rmqClient.MessageView, error) {
	mvs, err := d.Reader.Receive(d.Ctx, maxMessageNum, invisibleDuration)
	if err != nil {
		log.Errorf("RocketMqSubscribe Receive error:%v", err)
		return nil, err
	}
	return mvs, nil
}
