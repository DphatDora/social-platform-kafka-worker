package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"social-platform-kafka-worker/internal/model"
	"social-platform-kafka-worker/internal/repository"
	"social-platform-kafka-worker/internal/service"
	"social-platform-kafka-worker/package/constant"
)

type InterestScoreHandler struct {
	taskRepo             *repository.TaskRepository
	interestScoreService *service.InterestScoreService
}

func NewInterestScoreHandler(
	taskRepo *repository.TaskRepository,
	interestScoreService *service.InterestScoreService,
) *InterestScoreHandler {
	return &InterestScoreHandler{
		taskRepo:             taskRepo,
		interestScoreService: interestScoreService,
	}
}

// ProcessInterestScoreBotTasks fetches and processes unprocessed interest score bot tasks
func (h *InterestScoreHandler) ProcessInterestScoreBotTasks(ctx context.Context, batchSize int) error {
	// Fetch unprocessed tasks with action = "update_interest_score"
	tasks, err := h.taskRepo.FindUnprocessedByAction(constant.BOT_TASK_ACTION_UPDATE_INTEREST_SCORE, batchSize)
	if err != nil {
		log.Printf("[Error] Failed to fetch interest score bot tasks: %v", err)
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	log.Printf("[InterestScore] Found %d bot tasks to process", len(tasks))

	successCount := 0
	errorCount := 0

	// Process each task
	for _, task := range tasks {
		if err := h.ProcessSingleTask(&task); err != nil {
			log.Printf("[Error] Failed to process task ID=%d: %v", task.ID, err)
			errorCount++
			continue
		}

		// Mark task as executed
		now := time.Now()
		if err := h.taskRepo.MarkAsExecuted(task.ID, now); err != nil {
			log.Printf("[Error] Failed to mark task ID=%d as executed: %v", task.ID, err)
			errorCount++
			continue
		}

		successCount++
	}

	log.Printf("[InterestScore] Processed %d tasks: %d successful, %d errors", len(tasks), successCount, errorCount)
	return nil
}

// ProcessSingleTask processes a single bot task
func (h *InterestScoreHandler) ProcessSingleTask(task *model.BotTask) error {
	// Use the interest score service to process the update
	return h.interestScoreService.ProcessInterestScoreUpdate(json.RawMessage(task.Payload))
}
