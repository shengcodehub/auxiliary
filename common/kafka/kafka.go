package kafka

import (
	"github.com/gowins/dionysus/kafka"
	"github.com/gowins/dionysus/log"
	_kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

var producerClient *_kafka.Writer

func Setup() {

	w, err := kafka.NewWriter(viper.GetStringSlice("Kafka.Event.BootstrapServers"), viper.GetString("Kafka.Event.Topic"))
	if err != nil {
		log.Fatalf("open Producer kafka failed %s", err.Error())
	}

	producerClient = w
}

func GetProducer() *_kafka.Writer {
	return producerClient
}
