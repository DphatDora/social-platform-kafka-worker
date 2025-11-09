package repository

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"social-platform-kafka-worker/internal/model"
)

type TagPreferenceRepository struct {
	db *gorm.DB
}

func NewTagPreferenceRepository(db *gorm.DB) *TagPreferenceRepository {
	return &TagPreferenceRepository{db: db}
}

func (r *TagPreferenceRepository) UpsertTagPreferences(userID uint64, tags []string) error {
	if len(tags) == 0 {
		return nil // Nothing to update
	}

	// Check if record exists
	var existing model.UserTagPreference
	err := r.db.Where("user_id = ?", userID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		newPreference := &model.UserTagPreference{
			UserID:        userID,
			PreferredTags: tags,
		}
		if err := r.db.Create(newPreference).Error; err != nil {
			return fmt.Errorf("failed to create tag preferences: %w", err)
		}
		log.Printf("[TagPreference] Created for UserID=%d: %d tags", userID, len(tags))
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to query tag preferences: %w", err)
	}

	// Update existing record
	if err := r.db.Model(&model.UserTagPreference{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"preferred_tags": tags,
		}).Error; err != nil {
		return fmt.Errorf("failed to update tag preferences: %w", err)
	}

	log.Printf("[TagPreference] Updated for UserID=%d: %d tags", userID, len(tags))
	return nil
}

func (r *TagPreferenceRepository) FindByUser(userID uint64) (*model.UserTagPreference, error) {
	var preference model.UserTagPreference
	err := r.db.Where("user_id = ?", userID).First(&preference).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &preference, err
}

func (r *TagPreferenceRepository) GetActiveUserIDs(days int) ([]uint64, error) {
	var userIDs []uint64
	query := fmt.Sprintf(`
		SELECT DISTINCT user_id FROM (
			SELECT user_id FROM post_votes WHERE voted_at > NOW() - INTERVAL '%d days'
			UNION
			SELECT user_id FROM user_saved_posts WHERE is_followed = true
		) as active_users
	`, days)

	if err := r.db.Raw(query).Pluck("user_id", &userIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch active users: %w", err)
	}

	return userIDs, nil
}

func (r *TagPreferenceRepository) GetUserPreferredTags(userID uint64, limit int) ([]string, error) {
	var tags []string

	query := `
		SELECT DISTINCT unnest(p.tags) as tag
		FROM posts p
		INNER JOIN post_votes pv ON p.id = pv.post_id
		WHERE pv.user_id = ? AND pv.vote = true AND p.tags IS NOT NULL
		UNION
		SELECT DISTINCT unnest(p.tags) as tag
		FROM posts p
		INNER JOIN user_saved_posts usp ON p.id = usp.post_id
		WHERE usp.user_id = ? AND usp.is_followed = true AND p.tags IS NOT NULL
		LIMIT ?
	`

	if err := r.db.Raw(query, userID, userID, limit).Pluck("tag", &tags).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	return tags, nil
}
