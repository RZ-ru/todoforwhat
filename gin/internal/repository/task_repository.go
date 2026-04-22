package repository

import (
	"context"
	"todo/gin/internal/models"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id models.ID) (models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) (models.Task, error)
	DeleteTask(ctx context.Context, id models.ID) error
}
