package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Pratam-Kalligudda/Learn-RabbitMq/utils"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	utils.FailOnError(err, "failed to connect to amqp server")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "failed to create Channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topics",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "failed to declare queue")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FailOnError(err, "failed to declare queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topics", s)
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_topics",
			false,
			nil,
		)
		utils.FailOnError(err, "failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	utils.FailOnError(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
