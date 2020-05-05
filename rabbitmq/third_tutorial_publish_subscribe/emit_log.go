package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

// Implementing a pattern of publish/subscribe.
// Creating a log system. Where the producer will emit log messages
// Receiver will receive and print them
// Thus, published log messages will be broadcasted to all receivers
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

	// create an exchange
	// A producer never sends a message directly to queue, instead to an exchange
	// Fanout type exchange broadcasts a message to all the queues
	err = channel.ExchangeDeclare(
		"logs",
		"fanout", // type
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("failed to declare exchange:" + err.Error())
	}

	body := bodyFrom(os.Args)
	err = channel.Publish(
		"logs",
		"", // no queue name
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		panic("error publishing a message to the exchange:" + err.Error())
	}
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
