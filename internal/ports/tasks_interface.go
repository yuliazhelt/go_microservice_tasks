package ports

import (
	"tasks/internal/domain/models"
	"context"
)


type Tasks interface {
	CreateTask(ctx context.Context, author string, title string, description string, approves []*models.Approve) string
	CancelTask(ctx context.Context, author string, task *models.Task)
	CreateTaskRequest(ctx context.Context, task *models.Task)
	GetTasksByAuthor(ctx context.Context, taskId string) []*models.Task
	GetTaskById(ctx context.Context, taskId string) *models.Task
}
