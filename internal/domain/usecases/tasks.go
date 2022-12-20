package usecases

import (
	"tasks/internal/domain/models"
	"tasks/internal/ports"
	"tasks/pkg/infra/logger"
	"github.com/google/uuid"

	"context"
)

type Tasks struct {
	storage ports.TasksStorage
	m *MailSender
}

func New(l logger.Logger, storage ports.TasksStorage) (*Tasks, error) {
	return &Tasks{
		storage:storage,
		m:NewSender(l),
	}, nil
}


func (t *Tasks) CreateTask(ctx context.Context, author string, title string, description string, approves []*models.Approve) string {
	taskId := uuid.New().String()
	task := models.Task{Id : taskId, Author : author, Title : title, Description : description, Approves : approves}
	t.storage.AddTaskToStorage(ctx, &task)
	t.CreateTaskRequest(ctx, &task)
	return taskId
}

func (t *Tasks) CancelTask(ctx context.Context, author string, task *models.Task) {
	if (task.Author == author) {
		task.IsCancelled = true
	}
}

func (t *Tasks) CreateTaskRequest(ctx context.Context, task *models.Task) {
	for _, approve := range task.Approves {
		approve.Status = t.m.SendRequestMail(ctx, task, approve.Email)
		if approve.Status == "declined" {
			t.m.SendDeclineMail(ctx, task, task.Author, approve.Email)
			for _, approver := range task.Approves {
				t.m.SendDeclineMail(ctx, task, approver.Email, approve.Email)
			}
			return
		}
	}
	t.m.SendApproveMail(ctx, task, task.Author)
	for _, approver := range task.Approves {
		t.m.SendApproveMail(ctx, task, approver.Email)
	}
}


func (t *Tasks) GetTasksByAuthor(ctx context.Context, taskId string) []*models.Task {
	return t.storage.GetTasksByAuthor(ctx, taskId)
}

func (t *Tasks) GetTaskById(ctx context.Context, taskId string) *models.Task {
	return t.storage.GetTaskById(ctx, taskId)
}