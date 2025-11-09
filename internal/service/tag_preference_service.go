package service

import (
	"log"

	"social-platform-kafka-worker/internal/repository"
)

type TagPreferenceService struct {
	tagPrefRepo *repository.TagPreferenceRepository
}

func NewTagPreferenceService(repo *repository.TagPreferenceRepository) *TagPreferenceService {
	return &TagPreferenceService{
		tagPrefRepo: repo,
	}
}

func (s *TagPreferenceService) UpdateUserTagPreferences(userID uint64) error {
	// Get tags from posts that user has voted on or followed (limit 50)
	tags, err := s.tagPrefRepo.GetUserPreferredTags(userID, 50)
	if err != nil {
		log.Printf("[Error] Failed to fetch tags for UserID=%d: %v", userID, err)
		return err
	}

	if len(tags) == 0 {
		log.Printf("[Info] No tags found for UserID=%d, skipping update", userID)
		return nil
	}

	// Upsert tag preferences
	if err := s.tagPrefRepo.UpsertTagPreferences(userID, tags); err != nil {
		log.Printf("[Error] Failed to upsert tag preferences for UserID=%d: %v", userID, err)
		return err
	}

	log.Printf("✅ Tag preferences updated for UserID=%d: %d tags", userID, len(tags))
	return nil
}

func (s *TagPreferenceService) UpdateAllActiveUsers() error {
	// Get all users who have voted or followed posts in the last 30 days
	userIDs, err := s.tagPrefRepo.GetActiveUserIDs(30)
	if err != nil {
		log.Printf("[Error] Failed to fetch active users: %v", err)
		return err
	}

	log.Printf("[TagPreference] Found %d active users to update", len(userIDs))

	successCount := 0
	errorCount := 0

	for _, userID := range userIDs {
		if err := s.UpdateUserTagPreferences(userID); err != nil {
			errorCount++
			continue
		}
		successCount++
	}

	log.Printf("[TagPreference] Completed: %d successful, %d errors", successCount, errorCount)
	return nil
}
