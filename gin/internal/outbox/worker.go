package outbox

import (
	"context"
	"log"
	"time"
)

type Repository interface {
	GetUnprocessedEvents(ctx context.Context) ([]Event, error)
	MarkEventProcessed(ctx context.Context, id string) error
	MarkEventFailed(ctx context.Context, id string, errorMsg string) error
}

// обёртка для RabitMQ
type Publisher interface {
	Publish(ctx context.Context, e Event) error
}

type Worker struct {
	repo      Repository
	publisher Publisher
}

func NewWorker(repo Repository, publisher Publisher) *Worker {
	return &Worker{
		repo:      repo,
		publisher: publisher,
	}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		events, err := w.repo.GetUnprocessedEvents(ctx)
		if err != nil {
			log.Println("error fetching events:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, e := range events {

			// имитация отправки
			//log.Printf("send event: %s %s\n", e.EventType, string(e.Payload))

			if err = w.publisher.Publish(ctx, e); err != nil {
				if err = w.repo.MarkEventFailed(ctx, e.ID, err.Error()); err != nil {
					log.Println("error marking failed:", err)
				}
				continue
			}

			err = w.repo.MarkEventProcessed(ctx, e.ID)
			if err != nil {
				log.Println("error marking processed:", err)
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func send(e Event) error {
	log.Printf("send event: %s %s\n", e.EventType, string(e.Payload))
	return nil
}
