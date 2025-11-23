package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"social-platform-kafka-worker/internal/model"
	"social-platform-kafka-worker/internal/repository"
	"social-platform-kafka-worker/package/constant"
	"social-platform-kafka-worker/package/payload"
	"social-platform-kafka-worker/package/util"

	"gorm.io/gorm"
)

type KarmaService struct {
	userBadgeRepo *repository.UserBadgeRepository
	userRepo      *repository.UserRepository
}

func NewKarmaService(userBadgeRepo *repository.UserBadgeRepository, userRepo *repository.UserRepository) *KarmaService {
	return &KarmaService{
		userBadgeRepo: userBadgeRepo,
		userRepo:      userRepo,
	}
}

func (s *KarmaService) UpdateKarma(payloadBytes []byte) {
	var karmaPayload payload.UpdateUserKarmaPayload
	if err := json.Unmarshal(payloadBytes, &karmaPayload); err != nil {
		log.Printf("[Error] unmarshaling karma payload: %v", err)
		return
	}

	monthYear := util.FormatMonthYear(karmaPayload.UpdatedAt)

	actorKarma := s.getKarmaScoreForActor(karmaPayload.Action)
	if actorKarma != 0 {
		// Update monthly badge karma
		if err := s.processUserKarma(karmaPayload.UserId, actorKarma, monthYear); err != nil {
			log.Printf("[Error] processing karma for actor (user_id=%d): %v", karmaPayload.UserId, err)
		} else {
			log.Printf("Updated karma for actor (user_id=%d): %+d", karmaPayload.UserId, actorKarma)
		}

		// Update cumulative user karma
		if err := s.userRepo.UpdateKarma(karmaPayload.UserId, actorKarma); err != nil {
			log.Printf("[Error] updating cumulative karma for actor (user_id=%d): %v", karmaPayload.UserId, err)
		}
	}

	if karmaPayload.TargetId != nil {
		targetKarma := s.getKarmaScoreForTarget(karmaPayload.Action)
		if targetKarma != 0 {
			// Update monthly badge karma
			if err := s.processUserKarma(*karmaPayload.TargetId, targetKarma, monthYear); err != nil {
				log.Printf("[Error] processing karma for target (user_id=%d): %v", *karmaPayload.TargetId, err)
			} else {
				log.Printf("Updated karma for target (user_id=%d): %+d", *karmaPayload.TargetId, targetKarma)
			}

			// Update cumulative user karma
			if err := s.userRepo.UpdateKarma(*karmaPayload.TargetId, targetKarma); err != nil {
				log.Printf("[Error] updating cumulative karma for target (user_id=%d): %v", *karmaPayload.TargetId, err)
			}
		}
	}
}

func (s *KarmaService) getKarmaScoreForActor(action string) int {
	switch action {
	case constant.KARMA_ACTION_CREATE_POST:
		return constant.KARMA_SCORE_CREATE_POST
	case constant.KARMA_ACTION_CREATE_COMMENT:
		return constant.KARMA_SCORE_CREATE_COMMENT
	case constant.KARMA_ACTION_UPVOTE_POST, constant.KARMA_ACTION_DOWNVOTE_POST:
		return constant.KARMA_SCORE_VOTE_POST
	case constant.KARMA_ACTION_UPVOTE_COMMENT, constant.KARMA_ACTION_DOWNVOTE_COMMENT:
		return constant.KARMA_SCORE_VOTE_COMMENT
	default:
		return 0
	}
}

func (s *KarmaService) getKarmaScoreForTarget(action string) int {
	switch action {
	case constant.KARMA_ACTION_CREATE_COMMENT:
		return constant.KARMA_SCORE_GET_POST_COMMENT
	case constant.KARMA_ACTION_UPVOTE_POST:
		return constant.KARMA_SCORE_GET_POST_UPVOTE
	case constant.KARMA_ACTION_DOWNVOTE_POST:
		return constant.KARMA_SCORE_GET_POST_DOWNVOTE
	case constant.KARMA_ACTION_UPVOTE_COMMENT:
		return constant.KARMA_SCORE_GET_COMMENT_UPVOTE
	case constant.KARMA_ACTION_DOWNVOTE_COMMENT:
		return constant.KARMA_SCORE_GET_COMMENT_DOWNVOTE
	default:
		return 0
	}
}

func (s *KarmaService) processUserKarma(userID uint64, karmaChange int, monthYear string) error {
	existingBadge, err := s.userBadgeRepo.FindByUserAndMonth(userID, monthYear)

	if err == gorm.ErrRecordNotFound {
		newKarma := karmaChange
		if newKarma < 0 {
			newKarma = 0
		}

		badgeID := s.calculateBadgeLevel(uint64(newKarma))

		newBadge := &model.UserBadge{
			UserID:    userID,
			BadgeID:   badgeID,
			MonthYear: monthYear,
			Karma:     uint64(newKarma),
			AwardedAt: time.Now(),
		}

		if err := s.userBadgeRepo.Create(newBadge); err != nil {
			return fmt.Errorf("failed to create user badge: %w", err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to find user badge: %w", err)
	}

	// calculate new karma
	newKarma := int64(existingBadge.Karma) + int64(karmaChange)
	if newKarma < 0 {
		newKarma = 0
	}

	// calculate badge levels
	oldBadgeID := existingBadge.BadgeID
	newBadgeID := s.calculateBadgeLevel(uint64(newKarma))

	var badgeIDToUpdate uint64 = 0
	if newBadgeID > oldBadgeID {
		badgeIDToUpdate = newBadgeID
	}

	// update karma and badge
	if err := s.userBadgeRepo.UpdateKarmaAndBadge(userID, monthYear, karmaChange, badgeIDToUpdate); err != nil {
		return fmt.Errorf("failed to update user badge: %w", err)
	}

	return nil
}

func (s *KarmaService) calculateBadgeLevel(karma uint64) uint64 {
	if karma >= constant.BADGE_LEVEL_DIAMOND {
		return constant.BADGE_ID_DIAMOND
	}
	if karma >= constant.BADGE_LEVEL_PLATINUM {
		return constant.BADGE_ID_PLATINUM
	}
	if karma >= constant.BADGE_LEVEL_GOLD {
		return constant.BADGE_ID_GOLD
	}
	if karma >= constant.BADGE_LEVEL_SILVER {
		return constant.BADGE_ID_SILVER
	}
	if karma >= constant.BADGE_LEVEL_BRONZE {
		return constant.BADGE_ID_BRONZE
	}
	return constant.BADGE_ID_BRONZE
}
