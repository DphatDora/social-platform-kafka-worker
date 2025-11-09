package repository

import (
	"time"

	"gorm.io/gorm"
	"social-platform-kafka-worker/internal/model"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) FindDueTasks(now time.Time) ([]model.BotTask, error) {
	var tasks []model.BotTask
	err := r.db.Where("executed_at <= ?", now).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) Delete(task model.BotTask) error {
	return r.db.Delete(&task).Error
}

func (r *TaskRepository) FindUnprocessedByAction(action string, limit int) ([]model.BotTask, error) {
	var tasks []model.BotTask
	err := r.db.Where("action = ? AND executed_at IS NULL", action).
		Order("created_at ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) MarkAsExecuted(taskID uint, executedAt time.Time) error {
	return r.db.Model(&model.BotTask{}).
		Where("id = ?", taskID).
		Update("executed_at", executedAt).Error
}
