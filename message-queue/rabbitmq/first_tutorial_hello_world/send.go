package main

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

	// publish a message to queue
	body := "Hello World!"
	err = channel.Publish(
		"", // default exchange
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}
	log.Printf(" [x] Sent %s", body)
}
