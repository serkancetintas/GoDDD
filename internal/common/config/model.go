package config

type ApplicationConfig struct {
	Rabbit RabbitConfig `yaml:"rabbit"`
}

type RabbitConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	VirtualHost    string `yaml:"virtualHost"`
	ConnectionName string `yaml:"connectionName"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
}

type QueuesConfig struct {
	Order OrderQueueConfig `yaml:"order"`
}

type OrderQueueConfig struct {
	OrderCreated QueueConfig `yaml:"orderCreated"`
}

type QueueConfig struct {
	PrefetchCount int    `yaml:"prefetchCount"`
	ChannelCount  int    `yaml:"prefetchCount"`
	Exchange      string `yaml:"exchange"`
	ExchangeType  string `yaml:"exchangeType"`
	RoutingKey    string `yaml:"routingKey"`
	Queue         string `yaml:"queue"`
}

//MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

//Message is the amqp request to publish
type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationID string
	Priority      uint8
	Body          MessageBody
}
