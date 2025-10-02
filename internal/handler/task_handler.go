package handler

import (
	"context"
	"log"

	"social-platform-kafka-worker/internal/kafka"
	"social-platform-kafka-worker/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
	producer    *kafka.Producer
}

func NewTaskHandler(s *service.TaskService, p *kafka.Producer) *TaskHandler {
	return &TaskHandler{s, p}
}

func (h *TaskHandler) ProcessDueTasks(ctx context.Context) error {
	tasks, err := h.taskService.GetDueTasks(ctx)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if err := h.producer.SendMessage(ctx, t); err == nil {
			h.taskService.DeleteTask(t)
			log.Printf("✅ Task %d sent to Kafka and deleted", t.ID)
		} else {
			log.Printf("❌ Failed to send task %d: %v", t.ID, err)
		}
	}
	return nil
}
