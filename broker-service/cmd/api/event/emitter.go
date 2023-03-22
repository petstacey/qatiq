package event

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	Connection *amqp.Connection
}

func (e *Emitter) Setup() error {
	ch, err := e.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return declareExchange(ch)
}

func (e *Emitter) Push(event, severity string) error {
	ch, err := e.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	fmt.Println("Pushing event to channel")
	err = ch.PublishWithContext(
		context.Background(),
		"logs_topic",
		"log.INFO",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{Connection: conn}
	err := emitter.Setup()
	if err != nil {
		return Emitter{}, err
	}
	return emitter, nil
}
