package types

type Conf struct {
	BootstrapServers []string
	TopicConf        TopicConfMap
}

type TopicConfMap struct {
	Default TopicConf // 区分不同的topic
}

type TopicConf struct {
	Topic         string
	ConsumerGroup string
}
