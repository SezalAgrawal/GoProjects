package main

// Will keep this running to listen to new messages

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	url := "amqp://guest:guest@localhost:5672"

	// connect to RabbitMQ instance
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	defer connection.Close()

	// create a channel from the connection. Use channels to communicate with the queues rather than the connection itself.
	channel, err := connection.Channel()
	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}
	defer channel.Close()

	// create queue. A message is published to a queue
	// Queue is declared on both sender and receiver as we might start consumer before publisher
	// and we want to make sure queue exists before we start consuming messages from it.
	q, err := channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// create a consumer to consume messages from channel
	err = channel.Consume(
		q.Name,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("error registering a consumer:" + err.Error())
	}

	forever := make(chan bool)

	go func() {

	}
}
