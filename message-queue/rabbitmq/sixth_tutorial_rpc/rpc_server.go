package main

// A client communicates with a RPC server
// Server runs a function on a remote computer and wait for the result > Remote Procedure Call(RPC)
// Working of RPC:
// 1. Client creates an anonymous callback queue, where it awaits responses
// 2. CLient sends multiple requests in a declared queue, using the same callback queue
// 	(mentioned in reply_to) and different corrlation ids, unique to each request
// 3. Server processes the requests and sends responses to the callback queue (reply_to)
// 4. Client fetches the responses by identifying the correlation id

import (
	"log"
	"strconv"

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
		"rpc_queue",
		false, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// fair dispatch of messages
	// here prefetch count 1 means that the queue would not send more than 1
	// message to a worker. It will wait till it receives ack.
	// This ensures that heavy and light tasks are efficiently distributed
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic("error sending Qos: " + err.Error())
	}

	// create a consumer to consume messages from channel
	msgs, err := channel.Consume(
		q.Name,
		"",
		false, // auto-ack set to false
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
			n, err := strconv.Atoi(string(msg.Body))
			if err != nil {
				panic("error converting body to integer:" + err.Error())
			}

			log.Printf(" [.] fib(%d)", n)
			response := fib(n)

			err = channel.Publish(
				"",
				msg.ReplyTo, // routing key
				false,
				false,
				amqp.Publishing{
					// marks messages as persistent.
					// These persistent guarantees are not strong, because there is a time gap
					// when the node has accepted message and saved to disk
					// In that small time gap, we might loose the message
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(strconv.Itoa(response)),
					//  correlation id
					// to get response for each request, we can have a callback queue for each request: not efficient
					// efficient method is to create a correlation id per client. Now that there are many responses for the client,
					// to identify which response is for which request, we use correlation id
					CorrelationId: msg.CorrelationId,
				},
			)
			if err != nil {
				panic("error publishing a message to the queue:" + err.Error())
			}
			// For this auto-ack in channel is set to false
			// when a consumer consumes a message it sends an ack
			// back to the receiver, thereby deleting from the queue.
			// If did not receive ack, then the message gets delivered using
			// some other worker node
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
