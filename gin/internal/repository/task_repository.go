package repository

import (
	"context"
	"todo/gin/internal/models"
	"todo/gin/internal/outbox"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id models.ID) (models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) (models.Task, error)
	DeleteTask(ctx context.Context, id models.ID) error

	GetUnprocessedEvents(ctx context.Context) ([]outbox.Event, error)
	MarkEventProcessed(ctx context.Context, id string) error
	MarkEventFailed(ctx context.Context, id string, errorMsg string) error
}
