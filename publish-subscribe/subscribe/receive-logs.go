package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Pratam-Kalligudda/Learn-RabbitMq/utils"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "error while connecting to rabbitmq server")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "error while creating channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "error while creating exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FailOnError(err, "error while creating queue")

	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	utils.FailOnError(err, "error while binding exchange with queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CLRT+C")
	<-forever
}
