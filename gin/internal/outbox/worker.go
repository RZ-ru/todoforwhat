package outbox

import (
	"context"
	"log"
	"time"
)

type Repository interface {
	GetUnprocessedEvents(ctx context.Context) ([]Event, error)
	MarkEventProcessed(ctx context.Context, id string) error
}

type Worker struct {
	repo Repository
}

func NewWorker(repo Repository) *Worker {
	return &Worker{repo: repo}
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
			log.Printf("send event: %s %s\n", e.EventType, string(e.Payload))

			err := w.repo.MarkEventProcessed(ctx, e.ID)
			if err != nil {
				log.Println("error marking processed:", err)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
