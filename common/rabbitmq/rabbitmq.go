package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type Conf struct {
	Username string
	Password string
	Host     string
	Vhost    string
}

var producerChannel *amqp.Channel
var client *amqp.Connection

func SetUp(c Conf) {
	username := c.Username
	password := c.Password
	host := c.Host
	port := 5672
	vhost := c.Vhost

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		logx.Errorf("rabbitmq 创建连接失败: %v", err)
		return
	}
	client = conn
	OpenChannel()
}

func GetCh() *amqp.Channel {
	return producerChannel
}

func GetClient() *amqp.Connection {
	return client
}

func OpenChannel() {
	ch, err := GetClient().Channel()
	if err != nil {
		logx.Errorf("rabbitmq 打开 channel 错误: %v", err)
		return
	}
	producerChannel = ch
}

func QueueDeclare(queue string) (amqp.Queue, error) {
	q, err := GetCh().QueueDeclare(
		queue, // 队列名称
		true,  // 持久化
		false, // 自动删除
		false, // 独占
		false, // no-wait
		nil,   // 额外参数
	)
	if err != nil {
		if strings.Contains(err.Error(), "channel/connection is not open") {
			OpenChannel()
			return QueueDeclare(queue)
		}
		logx.Errorf("rabbitmq 声明队列失败: %v", err)
		return q, err
	}
	return q, nil
}
