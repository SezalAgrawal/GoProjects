package main

import (
	"log"
	"os"

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

	// declare a temporary queue
	// for loggining we need a fresh queue everytime we connect to rabbitMQ
	// server choses a random name, like amq.gen-JzTY20BRgKO-HjmUJj0wLg
	// when connection closes, queue will be deleted if exclusive
	q, err := channel.QueueDeclare(
		"",    // empty name
		false, // not durable
		false,
		true, // exclusive
		false,
		nil,
	)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// create a binding between exchange and queue
	// messages will be lost if no queue is bound to exchange
	// in this design if no receivers are started, messages will be discarded
	// the routing key binds the routing key of publish to a queue,
	// thereby segregating messages
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topic", s)
		err = channel.QueueBind(
			q.Name,
			s, // routing key
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			panic("failed to bind a queue: " + err.Error())
		}
	}

	// create a consumer to consume messages from channel
	msgs, err := channel.Consume(
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
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
