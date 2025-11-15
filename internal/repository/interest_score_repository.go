package repository

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"social-platform-kafka-worker/internal/model"
	"social-platform-kafka-worker/package/constant"
)

type InterestScoreRepository struct {
	db *gorm.DB
}

func NewInterestScoreRepository(db *gorm.DB) *InterestScoreRepository {
	return &InterestScoreRepository{db: db}
}

func (r *InterestScoreRepository) CreateOrUpdate(userID, communityID uint64, scoreDelta float64, action string) error {
	now := time.Now()

	// reset score
	if action == constant.INTEREST_ACTION_LEAVE_COMMUNITY {
		return r.handleLeaveCommunity(userID, communityID, now)
	}

	// Try to find existing record
	var existingScore model.UserInterestScore
	err := r.db.Where("user_id = ? AND community_id = ?", userID, communityID).
		First(&existingScore).Error

	if err == gorm.ErrRecordNotFound {
		newScore := &model.UserInterestScore{
			UserID:      userID,
			CommunityID: communityID,
			Score:       scoreDelta,
			UpdatedAt:   now,
		}

		// Set specific timestamp based on action
		switch action {
		case constant.INTEREST_ACTION_UPVOTE_POST, constant.INTEREST_ACTION_DOWNVOTE_POST:
			newScore.LastVoteAt = &now
		case constant.INTEREST_ACTION_JOIN_COMMUNITY:
			newScore.LastJoinAt = &now
		}

		if err := r.db.Create(newScore).Error; err != nil {
			return fmt.Errorf("failed to create interest score: %w", err)
		}

		log.Printf("[InterestScore] Created: UserID=%d, CommunityID=%d, Score=%.2f",
			userID, communityID, scoreDelta)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to query interest score: %w", err)
	}

	// Update existing record
	updates := map[string]interface{}{
		"score":      gorm.Expr("score + ?", scoreDelta),
		"updated_at": now,
	}

	// Update specific timestamps based on action
	switch action {
	case constant.INTEREST_ACTION_UPVOTE_POST, constant.INTEREST_ACTION_DOWNVOTE_POST:
		updates["last_vote_at"] = now
	case constant.INTEREST_ACTION_JOIN_COMMUNITY:
		updates["last_join_at"] = now
	}

	if err := r.db.Model(&model.UserInterestScore{}).
		Where("user_id = ? AND community_id = ?", userID, communityID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update interest score: %w", err)
	}

	newScore := existingScore.Score + scoreDelta
	log.Printf("[InterestScore] Updated: UserID=%d, CommunityID=%d, OldScore=%.2f, NewScore=%.2f",
		userID, communityID, existingScore.Score, newScore)

	return nil
}

func (r *InterestScoreRepository) FindByUser(userID uint64) ([]model.UserInterestScore, error) {
	var scores []model.UserInterestScore
	err := r.db.Where("user_id = ?", userID).
		Order("score DESC").
		Find(&scores).Error
	return scores, err
}

func (r *InterestScoreRepository) FindTopCommunitiesByUser(userID uint64, limit int) ([]model.UserInterestScore, error) {
	var scores []model.UserInterestScore
	err := r.db.Where("user_id = ? AND score > 0", userID).
		Order("score DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *InterestScoreRepository) handleLeaveCommunity(userID, communityID uint64, now time.Time) error {
	var existingScore model.UserInterestScore
	err := r.db.Where("user_id = ? AND community_id = ?", userID, communityID).
		First(&existingScore).Error

	if err == gorm.ErrRecordNotFound {
		newScore := &model.UserInterestScore{
			UserID:      userID,
			CommunityID: communityID,
			Score:       0.0,
			UpdatedAt:   now,
		}

		if err := r.db.Create(newScore).Error; err != nil {
			return fmt.Errorf("failed to create interest score record for leave_community: %w", err)
		}

		log.Printf("[InterestScore] Leave Community: Created new record with Score=0 for UserID=%d, CommunityID=%d",
			userID, communityID)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to query interest score for leave_community: %w", err)
	}

	// Reset existing score to 0
	updates := map[string]interface{}{
		"score":      0.0,
		"updated_at": now,
	}

	if err := r.db.Model(&model.UserInterestScore{}).
		Where("user_id = ? AND community_id = ?", userID, communityID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to reset interest score for leave_community: %w", err)
	}

	log.Printf("[InterestScore] Leave Community: Reset score to 0 for UserID=%d, CommunityID=%d (OldScore=%.2f)",
		userID, communityID, existingScore.Score)

	return nil
}
