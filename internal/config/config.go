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
	CompanyToken          = ""
	RabbitMQHost          = ""
	RabbitMQPort          = ""
	RabbitMQUser          = ""
	RabbitMQPassword      = ""
	RabbitMQVHost         = ""
	AmqpURI               = ""
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

	// Company Token
	CompanyToken = os.Getenv("COMPANY_TOKEN")

	// RabbitMQ Config
	RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	RabbitMQPort = os.Getenv("RABBITMQ_PORT")
	RabbitMQUser = os.Getenv("RABBITMQ_USER")
	RabbitMQPassword = os.Getenv("RABBITMQ_PASSWORD")
	RabbitMQVHost = os.Getenv("RABBITMQ_VHOST")

	AmqpURI = fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		RabbitMQUser,
		RabbitMQPassword,
		RabbitMQHost,
		RabbitMQPort,
		RabbitMQVHost,
	)

}
