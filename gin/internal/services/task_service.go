package services

import (
	"context"
	"time"
	"todo/gin/internal/models"
	"todo/gin/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, description string) (models.Task, error) {
	task := models.Task{
		Description: description,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateTaskWithEvent(ctx, task)
}

func (s *TaskService) GetTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetTasks(ctx)
}
func (s *TaskService) GetTaskByID(ctx context.Context, id models.ID) (models.Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}
func (s *TaskService) UpdateTask(ctx context.Context, id models.ID, description string) (models.Task, error) {
	task := models.Task{
		Id:          id,
		Description: description,
	}

	return s.repo.UpdateTask(ctx, task)
}
func (s *TaskService) DeleteTask(ctx context.Context, id models.ID) error {
	return s.repo.DeleteTask(ctx, id)
}
