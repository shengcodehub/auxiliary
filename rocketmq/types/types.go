package types

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
