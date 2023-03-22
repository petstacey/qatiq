package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = 80

type configuration struct {
	port int64
	env  string
	cors struct {
		trustedOrigins []string
	}
}

type broker struct {
	cfg    configuration
	rabbit *amqp.Connection
	wg     sync.WaitGroup
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var app broker
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	app.cfg.port = *flags.Int64("port", port, "port to listen on")
	app.cfg.env = *flags.String("env", "development", "environment [development | production]")
	flags.Func("trusted-origins", "trusted origins (space separated list)", func(val string) error {
		app.cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	rabbitConn, err := connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	app.rabbit = rabbitConn
	app.wg = sync.WaitGroup{}
	return app.serve()
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
		if counts > 5 {
			fmt.Println("Cannot connect to RabbitMQ. Quitting...")
			fmt.Println(err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		fmt.Println("Backing off")
		time.Sleep(backoff)
	}
	return connection, nil
}
