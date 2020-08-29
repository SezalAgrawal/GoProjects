package main

import (
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"golang.org/x/exp/errors/fmt"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Panicf("Could not instantiate pulsar client: %v", err)
	}
	defer client.Close()

	channel := make(chan pulsar.ConsumerMessage, 100)
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "my-topic",
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
		MessageChannel:   channel,
	})
	if err != nil {
		log.Panicf("Could not create consumer: %v", err)
	}
	defer consumer.Close()

	for ch := range channel {
		msg := ch.Message

		fmt.Printf("Received message msgId: %#v -- content: '%s' from consumer: %v\n",
			msg.ID(), string(msg.Payload()), ch.Consumer)

		consumer.Ack(msg)
	}

	if err := consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
}
