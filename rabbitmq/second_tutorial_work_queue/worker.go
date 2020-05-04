package main

// It fakes a second of work for every dot in message body.
// Pops message form queue and performs the task
// If started multiple workers, they work in a round-robin dispatching way
// Thus, if there is a lot of backlog of works, one could deploy extra worker nodes, 
// this will help in scaling. On average every consumer gets equal messages
import (
	"log"
	"time"
	"bytes"
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
		"task_queue",
		true, // durable -> This tells that queue will survive rabbitMQ node restart 
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// err = channel.Qos(
	// 	1, 		// prefetch count
	// 	0, 		// prefetch size
	// 	false, 	// global
	// )
	// if err != nil {
	// 	panic("error sending Qos: " + err.Error())
	// }

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
			log.Printf("Received a message: %s", msg.Body)
			dotCount := bytes.Count(msg.Body, []byte("."))
			duration := time.Duration(dotCount)
			time.Sleep(duration * time.Second)
			log.Printf("done")
			// For this auto-ack in channel is set to false
			// when a consumer consumes a message it sends an ack
			// back to the receiver, thereby deleting from the queue.
			// If did not receive ack, then the message gets delivered using
			// some other worker node
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
