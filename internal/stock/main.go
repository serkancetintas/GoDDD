package main

import (
	"fmt"
	"go-practice/internal/common/config"
	"go-practice/internal/common/rabbit"
	"go-practice/internal/stock/queue"
)

func main() {
	configurationManager := config.NewConfigurationManager()
	rabbitConfig := configurationManager.GetRabbitConfig()
	queuesConfig := configurationManager.GetQueuesConfig()

	rabbitClient := rabbit.NewRabbitClient(rabbitConfig, queuesConfig)
	defer rabbitClient.CloseConnection()

	consumerChan := make(chan bool)

	queue.InitializeConsumers(rabbitClient)
	fmt.Println("Started consumers")

	<-consumerChan // close(consumerChan)
}
