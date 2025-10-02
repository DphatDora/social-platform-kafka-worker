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
