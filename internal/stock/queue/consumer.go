package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"go-practice/internal/common/rabbit"
)

type Consumer struct {
}

func (c *Consumer) ConsumeOrderCreated(delivery amqp.Delivery) error {
	payload, _ := rabbit.Deserialize(delivery.Body)

	fmt.Printf("ConsumeOrderCreated, %s \n", payload)

	return nil
}
