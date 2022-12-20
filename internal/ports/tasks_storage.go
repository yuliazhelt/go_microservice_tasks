package ports

import (
	"tasks/internal/domain/models"
	"context"
)

type TasksStorage interface {
	AddTaskToStorage(ctx context.Context, task *models.Task)
	GetTasksByAuthor(ctx context.Context, author string) []*models.Task
	GetTaskById(ctx context.Context, taskId string) *models.Task
}