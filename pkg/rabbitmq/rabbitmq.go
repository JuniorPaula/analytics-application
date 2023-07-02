package rabbitmq

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/database"
	"c2d-reports/internal/repositories"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// ConsumerOnReportsQueue consumes the messages from the reports queue
func connect() *amqp.Connection {
	// connect to rabbitmq
	conn, err := amqp.Dial(config.AmqpURI)
	if err != nil {
		log.Fatalf("error while to connect rabbitMQ: %v", err)
	}

	return conn
}

// declareExchange declares an exchange of type direct
func declareExchange(ch *amqp.Channel) error {
	// declare exchange
	err := ch.ExchangeDeclare(
		config.ExchangeName, // name
		amqp.ExchangeDirect, // kind
		true,                // durable
		false,               // autoDelete
		false,               // internal
		false,               // noWait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("error while to declare exchange: %v", err)
	}

	return nil
}

// declareQueue declares a queue and binds it to the exchange
func declareQueue(ch *amqp.Channel) error {
	// declare a queue
	_, err := ch.QueueDeclare(
		config.QueueName, // name
		true,             // durable
		false,            // autoDelete
		false,            // internal
		false,            // noWait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("error while to declare a queue: %v", err)
	}

	// bind queue to exchange
	err = ch.QueueBind(
		config.QueueName,    // queue name
		"amqp.direct",       // routing key
		config.ExchangeName, // exchange
		false,               // noWait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("error while to bind queue to exchange: %v", err)
	}

	return nil
}

// PusblisherOnReportsQueue publishes a message on the reports queue
func PusblisherOnReportsQueue(message repositories.Report) {
	// connect to rabbitmq
	conn := connect()
	defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error while to create a channel: %v", err)
	}

	// declare exchange
	err = declareExchange(ch)
	if err != nil {
		log.Fatal(err)
	}

	// declare a queue
	err = declareQueue(ch)
	if err != nil {
		log.Fatal(err)
	}

	// json Marshal
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("error while to marshal message: %v", err)
	}

	// publish message to queue
	err = ch.Publish(
		config.ExchangeName, // exchange
		"amqp.direct",       // routing key
		false,               // mandatory
		false,               // immediate
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

// ConsumerOnReportsQueue consumes a message from the reports queue
func ConsumerOnReportsQueue() {
	// connect to rabbitmq
	conn := connect()
	defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error while to create a channel: %v", err)
	}

	// declare exchange
	err = declareExchange(ch)
	if err != nil {
		log.Fatal(err)
	}

	// declare a queue
	err = declareQueue(ch)
	if err != nil {
		log.Fatal(err)
	}

	// consume message from queue
	messages, err := ch.Consume(
		config.QueueName, // queue
		"",               // consumer
		false,            // autoAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // args
	)
	if err != nil {
		log.Fatalf("error while to consume message from queue: %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error while connect database;\n %s", err)
	}
	defer db.Close()

	repository := repositories.NewReportRepository(db)

	go func() {
		for reports := range messages {
			var report repositories.Report
			err := json.Unmarshal(reports.Body, &report)
			if err != nil {
				fmt.Printf("error while to unmarshal message: %v", err)
				reports.Nack(false, false)
				continue
			}

			// save report to database
			reportID, err := repository.CreateOrUpdate(report)
			if err != nil {
				fmt.Printf("ID: [%d]; report upserted:\n", reportID)

				reports.Nack(false, false)
				continue
			}
			fmt.Printf("ID: [%d]; new report computed:\n", reportID)

			reports.Ack(false)
		}
	}()

	fmt.Println("waiting for messages...")

	// block the main thread
	select {}
}
