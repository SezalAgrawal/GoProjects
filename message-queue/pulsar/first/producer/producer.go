package main

import (
	"context"
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

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "my-topic",
	})
	if err != nil {
		log.Panicf("Could not create producer: %v", err)
	}
	defer producer.Close()

	msgID, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload:      []byte("hello"),
	})
	if err != nil {
		log.Panicf("Failed to publish message: %v", err)
	}
	fmt.Printf("Published message: %v", msgID)
}
