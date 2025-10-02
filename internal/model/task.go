package model

import (
	"encoding/json"
	"time"
)

type BotTask struct {
	ID         uint            `gorm:"primaryKey"`
	Action     string          `gorm:"column:action"`
	Payload    json.RawMessage `gorm:"column:payload;type:jsonb"`
	CreatedAt  time.Time       `gorm:"column:created_at"`
	ExecutedAt time.Time       `gorm:"column:executed_at"`
}

func (BotTask) TableName() string {
	return "bot_tasks"
}
