package rocketmq

import (
	"auxiliary/rocketmq/types"
	"context"
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"time"
)

type MqSubscribe struct {
	reader rmqClient.SimpleConsumer
	closed int32
	conf   types.Conf
	ctx    context.Context
}

type Consumer interface {
	GetConsumer(ctx context.Context, c types.Conf) (*MqSubscribe, error)
	Subscribe(d *MqSubscribe) error
	Ack(d *MqSubscribe, mv *rmqClient.MessageView) error
	Receive(d *MqSubscribe, maxMessageNum int32, invisibleDuration time.Duration) ([]*rmqClient.MessageView, error)
	Stop(d *MqSubscribe)
}
