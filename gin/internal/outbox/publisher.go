package outbox

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	ch *amqp.Channel
}

func NewRabbitPublisher(ch *amqp.Channel) *RabbitPublisher {
	return &RabbitPublisher{ch: ch}
}

func (p *RabbitPublisher) Publish(ctx context.Context, e Event) error {
	return p.ch.PublishWithContext(
		ctx,
		"tasks_exchange",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        e.Payload,
			Type:        e.EventType,
		},
	)
}
