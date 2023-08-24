package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConnectStringDatabase = ""
	Port                  = 0
	ApiURL                = ""
	RabbitMQHost          = ""
	RabbitMQPort          = ""
	RabbitMQUser          = ""
	RabbitMQPassword      = ""
	RabbitMQVHost         = ""
	AmqpURI               = ""
	QueueName             = ""
	ExchangeName          = ""
)

func InitVariables() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		Port = 8081
	}

	// Database Config
	ConnectStringDatabase = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	)

	// API URL
	ApiURL = os.Getenv("API_URL")

	// RabbitMQ Config
	RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	RabbitMQPort = os.Getenv("RABBITMQ_PORT")
	RabbitMQUser = os.Getenv("RABBITMQ_USER")
	RabbitMQPassword = os.Getenv("RABBITMQ_PASSWORD")
	RabbitMQVHost = os.Getenv("RABBITMQ_VHOST")

	// Queue and Exchange Config
	QueueName = os.Getenv("QUEUE_NAME")
	ExchangeName = os.Getenv("EXCHANGE_NAME")

	// AmqpURI Config
	AmqpURI = fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		RabbitMQUser,
		RabbitMQPassword,
		RabbitMQHost,
		RabbitMQPort,
		RabbitMQVHost,
	)

}
