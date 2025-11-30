package repository

import (
	"time"

	"social-platform-kafka-worker/internal/model"

	"gorm.io/gorm"
)

type UserBadgeRepository struct {
	db *gorm.DB
}

func NewUserBadgeRepository(db *gorm.DB) *UserBadgeRepository {
	return &UserBadgeRepository{db: db}
}

func (r *UserBadgeRepository) FindByUserAndMonth(userID uint64, monthYear string) (*model.UserBadge, error) {
	var userBadge model.UserBadge
	err := r.db.Where("user_id = ? AND month_year = ?", userID, monthYear).
		First(&userBadge).Error
	if err != nil {
		return nil, err
	}
	return &userBadge, nil
}

func (r *UserBadgeRepository) Create(userBadge *model.UserBadge) error {
	return r.db.Create(userBadge).Error
}

func (r *UserBadgeRepository) UpdateKarmaAndBadge(userID uint64, monthYear string, karmaChange int, newBadgeID uint64) error {
	updates := map[string]interface{}{
		"karma": gorm.Expr("karma + ?", karmaChange),
	}

	// Nếu có badge_id mới (level up), cập nhật badge_id và awarded_at
	if newBadgeID > 0 {
		updates["badge_id"] = newBadgeID
		updates["awarded_at"] = time.Now()
	}

	return r.db.Model(&model.UserBadge{}).
		Where("user_id = ? AND month_year = ?", userID, monthYear).
		Updates(updates).
		Error
}

func (r *UserBadgeRepository) UpsertUserBadge(userID uint64, badgeID uint64, monthYear string, karma uint64) error {
	userBadge := &model.UserBadge{
		UserID:    userID,
		BadgeID:   badgeID,
		MonthYear: monthYear,
		Karma:     karma,
		AwardedAt: time.Now(),
	}

	return r.db.Exec(`
		INSERT INTO user_badges (user_id, badge_id, month_year, karma, awarded_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (user_id, badge_id, month_year)
		DO UPDATE SET
			karma = EXCLUDED.karma,
			awarded_at = EXCLUDED.awarded_at
	`, userBadge.UserID, userBadge.BadgeID, userBadge.MonthYear, userBadge.Karma, userBadge.AwardedAt).Error
}
