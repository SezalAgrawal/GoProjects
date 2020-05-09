package main

// Implements a work/task queue that distributes time consuming tasks among
// multiple workers. It avoids doing resource intensive task immediately and
// and schedules it afterwards.

import (
	"log"
	"os"
	"strings"

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
		"task_queue",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// publish a message to queue
	body := bodyFromArgs(os.Args)
	err = channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			// marks messages as persistent.
			// These persistent guarantees are not strong, because there is a time gap
			// when the node has accepted message and saved to disk
			// In that small time gap, we might loose the message
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}
	log.Printf(" [x] Sent %s", body)
}

func bodyFromArgs(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
