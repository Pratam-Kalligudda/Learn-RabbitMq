package main

import (
	"bytes"
	"log"
	"time"

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

	q, err := ch.QueueDeclare(
		"worker",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "error while creating queue")

	err = ch.Qos(1, 0, false)
	utils.FailOnError(err, "error while setting Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
