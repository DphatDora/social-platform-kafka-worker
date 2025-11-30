package repository

import (
	"social-platform-kafka-worker/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpdateKarma(userID uint64, karmaChange int) error {
	if karmaChange > 0 {
		return r.db.Model(&model.User{}).
			Where("id = ?", userID).
			Update("karma", gorm.Expr("karma + ?", karmaChange)).
			Error
	}

	return r.db.Exec(
		"UPDATE users SET karma = GREATEST(0, karma + ?) WHERE id = ?",
		karmaChange, userID,
	).Error
}
