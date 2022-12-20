package usecases

import (
	"tasks/internal/domain/models"

	"context"
)

type TasksStorage struct {
	storageByUser map[string][]*models.Task
	storageById map[string]*models.Task
}

func NewStorage() *TasksStorage {
	storageByUser := make(map[string][]*models.Task)
	storageById := make(map[string]*models.Task)

	return &TasksStorage{
		storageByUser: storageByUser,
		storageById : storageById,
	}
}

func (s *TasksStorage) AddTaskToStorage(ctx context.Context, task *models.Task) {
	s.storageByUser[task.Author] = append(s.storageByUser[task.Author], task)
	s.storageById[task.Id] = task
}


func (s *TasksStorage) GetTasksByAuthor(ctx context.Context, author string) []*models.Task {
	tasks, prs := s.storageByUser[author]
	if !prs {
		return []*models.Task{}
	}
	return tasks
}

func (s *TasksStorage) GetTaskById(ctx context.Context, taskId string) *models.Task {
	task, prs := s.storageById[taskId]
	if !prs {
		return nil
	}
	return task
}