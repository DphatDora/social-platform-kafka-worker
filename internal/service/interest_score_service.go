package service

import (
	"encoding/json"
	"fmt"
	"log"

	"social-platform-kafka-worker/internal/repository"
	"social-platform-kafka-worker/package/constant"
	"social-platform-kafka-worker/package/payload"
)

type InterestScoreService struct {
	interestScoreRepo *repository.InterestScoreRepository
}

func NewInterestScoreService(repo *repository.InterestScoreRepository) *InterestScoreService {
	return &InterestScoreService{
		interestScoreRepo: repo,
	}
}

func (s *InterestScoreService) ProcessInterestScoreUpdate(payloadData json.RawMessage) error {
	var p payload.UpdateInterestScorePayload
	if err := json.Unmarshal(payloadData, &p); err != nil {
		log.Printf("[Error] Failed to unmarshal interest score payload: %v", err)
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Validate payload
	if p.UserID == 0 || p.CommunityID == 0 {
		log.Printf("[Warning] Invalid payload: UserID=%d, CommunityID=%d", p.UserID, p.CommunityID)
		return fmt.Errorf("invalid payload: userID or communityID is zero")
	}

	// Calculate score delta based on action
	scoreDelta := s.GetScoreDeltaForAction(p.Action)

	// log.Printf("[InterestScore] Processing: UserID=%d, CommunityID=%d, Action=%s, ScoreDelta=%.2f",
	// 	p.UserID, p.CommunityID, p.Action, scoreDelta)

	// Update interest score
	if err := s.interestScoreRepo.CreateOrUpdate(p.UserID, p.CommunityID, scoreDelta, p.Action); err != nil {
		log.Printf("[Error] Failed to update interest score: %v", err)
		return fmt.Errorf("failed to update interest score: %w", err)
	}

	log.Printf("✅ Interest score updated successfully for UserID=%d, CommunityID=%d", p.UserID, p.CommunityID)
	return nil
}

func (s *InterestScoreService) GetScoreDeltaForAction(action string) float64 {
	switch action {
	case constant.INTEREST_ACTION_UPVOTE_POST:
		return constant.SCORE_UPVOTE_POST // +2.0
	case constant.INTEREST_ACTION_DOWNVOTE_POST:
		return constant.SCORE_DOWNVOTE_POST // -1.0
	case constant.INTEREST_ACTION_FOLLOW_POST:
		return constant.SCORE_FOLLOW_POST // +3.0
	case constant.INTEREST_ACTION_JOIN_COMMUNITY:
		return constant.SCORE_JOIN_COMMUNITY // +10.0
	case constant.INTEREST_ACTION_LEAVE_COMMUNITY:
		return constant.SCORE_LEAVE_COMMUNITY // Reset score to 0
	default:
		log.Printf("[Warning] Unknown action: %s, returning 0", action)
		return 0.0
	}
}
