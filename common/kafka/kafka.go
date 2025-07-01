package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	kClient "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type TaskConf struct {
	BootstrapServers []string
	SaslUsername     string
	SaslPassword     string
	SecurityProtocol string
	TopicConf        TopicTaskConfMap
}

type TopicTaskConfMap struct {
	TaskEvent   TopicConf // 区分不同的topic
	RewardEvent TopicConf
}

type Conf struct {
	BootstrapServers []string
	SaslUsername     string
	SaslPassword     string
	SecurityProtocol string
	TopicConf        TopicConfMap
}

type TopicConfMap struct {
	Log    TopicConf // 区分不同的topic
	User5E TopicConf
	Game   TopicConf
}

type TopicConf struct {
	Topic   string
	GroupId string
}

type MqBody struct {
	Event string `json:"event"` // 事件 用于区分走不同的逻辑
	Data  string `json:"data"`  // 数据
}

var producerClient *kClient.Writer

func SetUp(c *Conf) {
	producerClient = doInitProducer(c)
}

func GetProduct() *kClient.Writer {
	return producerClient
}

func GetLogProducer(ctx context.Context, c Conf, body *MqBody, key string) error {
	msg, err := jsoniter.Marshal(body)
	if err != nil {
		logx.Errorf("kafka GetLogProducer jsoniter MarshalToString  err: %v", err)
		return err
	}
	topic := c.TopicConf.Log.Topic
	err = GetProduct().WriteMessages(ctx, kClient.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: msg,
	})
	if err != nil {
		logx.Errorf("kafka GetLogProducer failed: %v", err)
		return err
	}
	return nil
}

func GetLogConsumer(c Conf) *kClient.Reader {
	dialer := doInitConsumer(&c)
	r := kClient.NewReader(kClient.ReaderConfig{
		Brokers:     c.BootstrapServers,
		GroupID:     c.TopicConf.Log.GroupId,
		Topic:       c.TopicConf.Log.Topic,
		Dialer:      dialer,
		Logger:      kClient.LoggerFunc(logf),
		ErrorLogger: kClient.LoggerFunc(logf),
	})
	fmt.Print("init kafka consumer success\n")
	return r
}

func doInitConsumer(cfg *Conf) *kClient.Dialer {
	fmt.Print("init kafka consumer, it may take a few seconds to init the connection\n")
	dialer := &kClient.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	switch cfg.SecurityProtocol {
	case "plaintext":

	case "sasl_ssl":
		mechanism := plain.Mechanism{
			Username: cfg.SaslUsername,
			Password: cfg.SaslPassword,
		}
		dialer.SASLMechanism = mechanism
		dialer.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	default:
		panic(errors.New("kafka unknown protocol"))
	}
	return dialer
}

func logf(msg string, a ...interface{}) {
	logx.Infof(msg, a...)
}

func doInitProducer(cfg *Conf) *kClient.Writer {
	fmt.Print("init kafka producer, it may take a few seconds to init the connection\n")
	w := &kClient.Writer{
		Addr:         kClient.TCP(cfg.BootstrapServers...),
		Balancer:     &kClient.Hash{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: 0,
		Logger:       kClient.LoggerFunc(logf),
		ErrorLogger:  kClient.LoggerFunc(logf),
		BatchSize:    1,
	}

	switch cfg.SecurityProtocol {
	case "plaintext":

	case "sasl_ssl":
		sharedTransport := &kClient.Transport{
			SASL: plain.Mechanism{
				Username: cfg.SaslUsername,
				Password: cfg.SaslPassword,
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		w.Transport = sharedTransport
	default:
		panic(errors.New("kafka unknown protocol"))
	}
	fmt.Print("init kafka producer success\n")
	return w
}

func doInitTaskProducer(cfg *TaskConf) *kClient.Writer {
	fmt.Print("init task kafka producer, it may take a few seconds to init the connection\n")
	w := &kClient.Writer{
		Addr:         kClient.TCP(cfg.BootstrapServers...),
		Balancer:     &kClient.Hash{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: 1,
		Logger:       kClient.LoggerFunc(logf),
		ErrorLogger:  kClient.LoggerFunc(logf),
		BatchSize:    1,
	}

	switch cfg.SecurityProtocol {
	case "plaintext":

	case "sasl_ssl":
		sharedTransport := &kClient.Transport{
			SASL: plain.Mechanism{
				Username: cfg.SaslUsername,
				Password: cfg.SaslPassword,
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		w.Transport = sharedTransport
	default:
		panic(errors.New("task kafka unknown protocol"))
	}
	fmt.Print("init task kafka producer success\n")
	return w
}
