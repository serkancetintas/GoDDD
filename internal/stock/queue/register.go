package queue

import (
	"errors"
	"github.com/streadway/amqp"
	"go-practice/internal/common/config"
	"go-practice/internal/common/rabbit"
)

type QueueConsumerMap map[config.QueueConfig]func(delivery amqp.Delivery) error

var qcm QueueConsumerMap

func getRegisteredQueueConsumer(c *rabbit.Client) QueueConsumerMap {
	if qcm != nil {
		return qcm
	}

	queueConsumerMap := make(QueueConsumerMap)

	// Order Queue-Consumer binding
	queueConsumerMap[c.QueuesConfig.Order.OrderCreated] = (&Consumer{}).ConsumeOrderCreated

	qcm = queueConsumerMap
	return qcm
}

func FindConsumer(routingKey string) (func(delivery amqp.Delivery) error, error) {
	for key, value := range qcm {
		if key.RoutingKey == routingKey {
			return value, nil
		}
	}
	return nil, errors.New("Consumer not found, Routing Key: " + routingKey)
}
