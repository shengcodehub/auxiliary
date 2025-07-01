package kafka

import (
	"context"
	"errors"
	dioKafka "github.com/gowins/dionysus/kafka"
	"github.com/gowins/dionysus/log"
	"github.com/segmentio/kafka-go"
	"github.com/shengcodehub/auxiliary/kafka/types"
	"time"
)

type Subscribe struct {
	Reader *kafka.Reader
	Closed int32
	Ctx    context.Context
	Conf   types.Conf
}

func GetConsumer(ctx context.Context, c types.Conf) (*Subscribe, error) {
	var (
		d                   = &Subscribe{}
		err                 error
		url, topic, groupID = c.BootstrapServers,
			c.TopicConf.Default.Topic,
			c.TopicConf.Default.ConsumerGroup
	)

	dialer := &kafka.Dialer{
		Timeout:   time.Second * 10,
		DualStack: true,
		KeepAlive: time.Second * 3,
	}

	d.Reader, err = dioKafka.NewGroupReader(
		url, topic, groupID,
		dioKafka.ReaderWithDialer(dialer),
		ReaderWithCommitInterval(time.Second),
		ReaderWithSessionTimeout(time.Second*10),
	)
	if err != nil {
		log.Errorf("new group reader failed: %v", err)
		return nil, err
	}
	d.Ctx = ctx
	return d, nil
}

func ReaderWithCommitInterval(t time.Duration) dioKafka.ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.CommitInterval = t
	}
}

func ReaderWithSessionTimeout(t time.Duration) dioKafka.ReaderOption {
	return func(c *kafka.ReaderConfig) {
		c.SessionTimeout = t
		c.RebalanceTimeout = t
	}
}

func Receive(d *Subscribe, invisibleDuration time.Duration) (kafka.Message, error) {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, invisibleDuration)
	defer cancel()
	msg, err := d.Reader.ReadMessage(timeout)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		log.Errorf("kafka read msg failed: %s", err.Error())
		return msg, err
	}
	if msg.Topic == "" {
		return msg, errors.New("topic is empty")
	}
	return msg, nil
}
