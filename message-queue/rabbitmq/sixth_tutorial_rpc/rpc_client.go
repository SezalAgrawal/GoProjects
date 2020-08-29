package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (resp int, err error) {
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
		"",
		false, // durable
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
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

	corrID := randomString(32)

	err = channel.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.Itoa(n)),
			//  correlation id
			// to get response for each request, we can have a callback queue for each request: not efficient
			// efficient method is to create a correlation id per client. Now that there are many responses for the client,
			// to identify which response is for which request, we use correlation id
			CorrelationId: corrID,
			ReplyTo: q.Name,
		},
	)
	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	for msg := range msgs {
		if corrID == msg.CorrelationId {
			resp, err = strconv.Atoi(string(msg.Body))
			if err != nil {
				panic("error converting body to int:" + err.Error())
			}
			break
		}
	}
	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	n := bodyFrom(os.Args)
	log.Printf(" [x] Requesting fib(%d)", n)
	resp, err := fibonacciRPC(n)
	if err != nil {
		panic("failed to handle RPC request:" + err.Error())
	}
	log.Printf(" [.] Got %d", resp)
}

func bodyFrom(args []string) int {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		panic("error converting arg to int:" + err.Error())
	}
	return n
}
