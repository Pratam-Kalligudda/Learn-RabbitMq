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
	defer conn.Channel()

	err = ch.ExchangeDeclare(
		"logs-direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "error while declaring exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := utils.BodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"logs-direct",
		utils.SeverityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	utils.FailOnError(err, "error while publishing content")

	log.Printf(" [x] Sent %s", body)
}
