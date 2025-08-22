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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	utils.FailOnError(err, "failed to connect to amqp server")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "failed to create channel")
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
	utils.FailOnError(err, "failed to declare exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := utils.BodyFrom(os.Args)

	err = ch.PublishWithContext(ctx,
		"logs_topics",
		utils.SeverityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	utils.FailOnError(err, "failed to publish content")

	log.Printf(" [x] Sent %s", body)
}
