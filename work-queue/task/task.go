package main

import (
	"context"
	"log"
	"os"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := utils.BodyFrom(os.Args)

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         []byte(body),
		},
	)
	utils.FailOnError(err, "error while publishing body")

	log.Printf(" [x] Send %s", body)
}
