package consumer

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch *amqp.Channel
}

func NewConsumer(ch *amqp.Channel) *Consumer {
	return &Consumer{ch: ch}
}

func (c *Consumer) Start(ctx context.Context) error {

	msgs, err := c.ch.Consume(
		"tasks", // очередь
		"",      // consumer name
		false,   // auto-ack (пока так)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case msg := <-msgs:
			log.Printf("received: %s\n", string(msg.Body))

			// имитация обработки
			err := process(msg.Body)

			if err != nil {
				log.Println("process error:", err)

				// вернуть сообщение обратно в очередь
				msg.Nack(false, true)
				continue
			}

			// подтвердить обработку
			msg.Ack(false)
		}
	}
}

func process(body []byte) error {
	log.Println("processing:", string(body))
	return nil
}
