package queue

import (
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go-practice/internal/common/rabbit"
	"log"
	"os"
)

func InitializeConsumers(c *rabbit.Client) {
	registeredQueueConsumers := getRegisteredQueueConsumer(c)
	for queueConfig := range registeredQueueConsumers {
		for i := 0; i < queueConfig.ChannelCount; i++ {
			channel := c.Channel
			channelName := generateChannelName()
			fmt.Println("consumer queue name:" + queueConfig.Queue)
			deliveries, err := channel.Consume(queueConfig.Queue, channelName, false, false, false, false, nil)
			if err != nil {
				log.Panicf("Terminating. Error details: %s", err.Error())
			}
			for j := 0; j < queueConfig.PrefetchCount; j++ {
				go func() {
					for d := range deliveries {
						consumer, err := FindConsumer(d.RoutingKey)
						if err != nil {
							logrus.Errorf("Consumer not found, error: %s", err.Error())
							nackOnError(d)
							continue
						}
						if err = consumer(d); err != nil {
							logrus.Errorf("An error occurred when consuming %s, error: %s", d.RoutingKey, err.Error())
							nackOnError(d)
							continue
						}
						ackMessage(d)
					}
				}()
			}
		}
	}
}

func generateChannelName() string {
	name := os.Getenv("HOSTNAME")
	if name == "" {
		return fmt.Sprintf("go-rabbit-consumer-app-%s", random.String(10, "123456789"))
	}
	return fmt.Sprintf("%s-%s", name, random.String(10, "123456789"))
}

func nackOnError(m amqp.Delivery) {
	if err := m.Nack(false, false); err != nil {
		logrus.Errorf("could not nack message %s %s", m.Body, err)
	}
}

func ackMessage(m amqp.Delivery) {
	err := m.Ack(false)
	if err != nil {
		logrus.Errorf("failed to ack message %s \t %s", m.Body, err)
	}
}
