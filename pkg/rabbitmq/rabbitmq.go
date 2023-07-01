package rabbitmq

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/repositories"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var queueName = "chat_reports"
var exchangeName = "reports"

func PusblisherOnReportsQueue(message repositories.Report) {
	// connect to rabbitmq
	conn, err := amqp.Dial(config.AmqpURI)
	if err != nil {
		log.Fatalf("error while to connect rabbitMQ: %v", err)
	}
	defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error while to create a channel: %v", err)
	}

	// declare exchange
	err = ch.ExchangeDeclare(
		exchangeName,        // name
		amqp.ExchangeDirect, // kind
		true,                // durable
		false,               // autoDelete
		false,               // internal
		false,               // noWait
		nil,                 // args
	)
	if err != nil {
		log.Fatalf("error while to declare exchange: %v", err)
	}

	// declare a queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // internal
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("error while to declare a queue: %v", err)
	}

	// bind queue to exchange
	err = ch.QueueBind(
		queueName,     // queue name
		"amqp.direct", // routing key
		exchangeName,  // exchange
		false,         // noWait
		nil,           // args
	)
	if err != nil {
		log.Fatalf("error while to bind queue to exchange: %v", err)
	}

	// json Marshal
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("error while to marshal message: %v", err)
	}

	// publish message to queue
	err = ch.Publish(
		exchangeName,  // exchange
		"amqp.direct", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		},
	)
	if err != nil {
		log.Fatalf("error while to publish message to queue: %v", err)
	}

	fmt.Println("message published to queue")
}
