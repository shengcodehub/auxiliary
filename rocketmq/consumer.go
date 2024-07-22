package rocketmq

import (
	"context"
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/shengwenjin/auxiliary/rocketmq/types"
	"time"
)

type MqSubscribe struct {
	Reader rmqClient.SimpleConsumer
	Closed int32
	Conf   types.Conf
	Ctx    context.Context
}

type Consumer interface {
	GetConsumer(ctx context.Context, c types.Conf) (*MqSubscribe, error)
	Subscribe(d *MqSubscribe) error
	Ack(d *MqSubscribe, mv *rmqClient.MessageView) error
	Receive(d *MqSubscribe, maxMessageNum int32, invisibleDuration time.Duration) ([]*rmqClient.MessageView, error)
	Stop(d *MqSubscribe)
}
