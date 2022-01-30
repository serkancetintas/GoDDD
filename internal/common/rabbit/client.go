package rabbit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go-practice/internal/common/config"
	"log"
	"strings"
	"time"
)

type Client struct {
	connection   *amqp.Connection
	queuesConfig config.QueuesConfig
	channel      *amqp.Channel
}

func NewRabbitClient(rabbitConfig config.RabbitConfig, queuesConfig config.QueuesConfig) *Client {
	c := createConnection(rabbitConfig)
	channel, err := c.Channel()
	if err != nil {
		channel.Close()
		log.Panicf("Channel could not created. Terminating. Error details: %s", err.Error())
	}

	return &Client{
		connection:   c,
		queuesConfig: queuesConfig,
		channel:      channel,
	}
}

func (c *Client) DeclareExchangeQueueBindings() {
	configs := c.getRegisteredQueue()

	for _, queueConfig := range configs {
		declareExchange(c.channel, queueConfig)
		declareQueue(c.channel, queueConfig)
		declareDeadLetterQueue(c.channel, queueConfig)
		bindQueue(c.channel, queueConfig)
		err := c.channel.Qos(queueConfig.PrefetchCount, 0, false)
		if err != nil {
			log.Panicf("PrefetchCount could not defined. Terminating. Error details: %s", err.Error())
		}
	}
}

func (c *Client) CreateChannel(prefetchCount int) *amqp.Channel {
	channel, err := c.connection.Channel()
	if err != nil {
		channel.Close()
		log.Panicf("Channel could not created. Terminating. Error details: %s", err.Error())
	}
	e := channel.Qos(prefetchCount, 0, false)
	if e != nil {
		log.Panicf("PrefetchCount could not defined. Terminating. Error details: %s", e.Error())
	}
	return channel
}

func declareExchange(channel *amqp.Channel, queueConfig config.QueueConfig) {
	err := channel.ExchangeDeclare(queueConfig.Exchange, queueConfig.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Exchange could not declared. Terminating. Error details: %s", err.Error())
	}
}

func declareQueue(channel *amqp.Channel, queueConfig config.QueueConfig) {
	deadLetterArgs := getDeadLetterArgs(queueConfig.Queue)
	_, err := channel.QueueDeclare(queueConfig.Queue, true, false, false, false, deadLetterArgs)
	if err != nil {
		log.Panicf("Queue could not declared. Terminating. Error details: %s", err.Error())
	}
}

func declareDeadLetterQueue(channel *amqp.Channel, queueConfig config.QueueConfig) {
	_, err := channel.QueueDeclare(queueConfig.Queue+".deadLetter", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Queue could not declared. Terminating. Error details: %s", err.Error())
	}
}

func bindQueue(channel *amqp.Channel, queueConfig config.QueueConfig) {
	err := channel.QueueBind(queueConfig.Queue, queueConfig.RoutingKey, queueConfig.Exchange, false, nil)
	if err != nil {
		log.Panicf("Binding could not defined. Terminating. Error details: %s", err.Error())
	}
}

func getDeadLetterArgs(queueName string) amqp.Table {
	return amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": queueName + ".deadLetter",
	}
}

func createConnection(rabbitConfig config.RabbitConfig) *amqp.Connection {
	amqpConfig := amqp.Config{
		Properties: amqp.Table{
			"connection_name": rabbitConfig.ConnectionName,
		},
		Heartbeat: 30 * time.Second,
	}
	connectionUrl := getConnectionUrl(rabbitConfig)
	connection, err := amqp.DialConfig(connectionUrl, amqpConfig)
	if err != nil {
		_ = connection.Close()
		log.Panicf("Client cannogt deserialize. Terminating. Error details: %s", err.Error())
	}
	log.Printf("RabbitMQ connected. Host: %s, Port: %d, Virtual Host: %s", rabbitConfig.Host, rabbitConfig.Port, rabbitConfig.VirtualHost)
	return connection
}

func (c *Client) Publish(m config.Message) error {
	p := amqp.Publishing{
		Headers:       amqp.Table{"type": m.Body.Type},
		ContentType:   m.ContentType,
		CorrelationId: m.CorrelationID,
		Body:          m.Body.Data,
		ReplyTo:       m.ReplyTo,
	}

	exchangeFirstPart := strings.Split(m.Queue, ".")[0]
	exchangeName := fmt.Sprintf("%s.%s", exchangeFirstPart, "events")

	if err := c.channel.Publish(exchangeName, m.Queue, false, false, p); err != nil {
		return fmt.Errorf("Error in Publishing: %s", err)
	}
	return nil
}

func (c *Client) CloseConnection() {
	c.connection.Close()
}

func getConnectionUrl(config config.RabbitConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.VirtualHost)
}

func Serialize(msg interface{}) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}

func Deserialize(b []byte) (interface{}, error) {
	var msg config.Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
