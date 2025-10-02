package service

import (
	"context"
	"time"

	"social-platform-kafka-worker/internal/model"
	"social-platform-kafka-worker/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo}
}

func (s *TaskService) GetDueTasks(ctx context.Context) ([]model.BotTask, error) {
	return s.repo.FindDueTasks(time.Now())
}

func (s *TaskService) DeleteTask(task model.BotTask) error {
	return s.repo.Delete(task)
}
