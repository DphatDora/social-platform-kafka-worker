package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserTagPreference struct {
	UserID        uint64         `gorm:"column:user_id;primaryKey"`
	PreferredTags pq.StringArray `gorm:"column:preferred_tags;type:text[]"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserTagPreference) TableName() string {
	return "user_tag_preferences"
}

func (u *UserTagPreference) BeforeCreate(tx *gorm.DB) error {
	if u.PreferredTags == nil {
		u.PreferredTags = pq.StringArray{}
	}
	return nil
}
