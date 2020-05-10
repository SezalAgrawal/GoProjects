package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

// Usecase: extension of direct logging
// We want to control the severity of the message and the source
// Example: we may want to listen to just critical errors coming from 'cron' but also all logs from 'kern'.
// * (star) can substitute for exactly one word.
// # (hash) can substitute for zero or more words.
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
	// Direct type of exchange to control the broadcasting of the messages
	err = channel.ExchangeDeclare(
		"logs_topic",
		"topic", // type
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("failed to declare exchange:" + err.Error())
	}

	// 'severity' can be one of 'info', 'warning', 'error'.
	body := bodyFrom(os.Args)
	err = channel.Publish(
		"logs_topic",
		severityFrom(os.Args), // routing key which decides which queue(s) to send message to
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
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
