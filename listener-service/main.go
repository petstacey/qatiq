package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/pso-dev/qatiq/backend/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	fmt.Println("Listening for and consuming messages")
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready..")
			counts++
		} else {
			fmt.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		if counts > 10 {
			fmt.Println(err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		fmt.Println("Backing off")
		time.Sleep(backoff)
	}
	return connection, nil
}
