package main

import (
	"log"
	"os"

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
		"logs-direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FailOnError(err, "error while creating queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,
			s,
			"logs-direct",
			false,
			nil,
		)
		utils.FailOnError(err, "error while binding queue to exchange")
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

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
