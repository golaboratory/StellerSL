package gamification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

type BadgeEngine struct {
	queries *db.Queries
}

func NewBadgeEngine(queries *db.Queries) *BadgeEngine {
	return &BadgeEngine{queries: queries}
}

func (e *BadgeEngine) AwardBadges(ctx context.Context, userID uuid.UUID, trigger string) error {
	// 1. Get current badges to avoid duplicates
	currentBadges, err := e.queries.ListUserBadges(ctx, userID)
	if err != nil {
		return err
	}

	badgeMap := make(map[string]bool)
	for _, b := range currentBadges {
		badgeMap[b.RequirementType] = true
	}

	// 2. Get stats for checking conditions
	stats, err := e.queries.GetDashboardStats(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return err
	}

	// 3. Logic for each badge
	allBadges, err := e.queries.ListAllBadges(ctx) // Need to add this query
	if err != nil {
		return err
	}

	for _, b := range allBadges {
		if badgeMap[b.RequirementType] {
			continue // Already have it
		}

		shouldAward := false
		switch b.RequirementType {
		case "first_task":
			if stats.CompletedTasks >= 1 {
				shouldAward = true
			}
		case "three_strike":
			if stats.CompletedTasks >= 3 {
				shouldAward = true
			}
		case "task_crusher_1":
			if stats.CompletedTasks >= 10 {
				shouldAward = true
			}
		case "early_bird":
			if trigger == "task_completed" {
				hour := time.Now().Hour()
				if hour >= 4 && hour < 9 {
					shouldAward = true
				}
			}
		case "midnight_owl":
			if trigger == "task_completed" {
				hour := time.Now().Hour()
				if hour >= 0 && hour < 4 {
					shouldAward = true
				}
			}
		}

		if shouldAward {
			_ = e.queries.AwardBadge(ctx, db.AwardBadgeParams{
				UserID:  userID,
				BadgeID: b.ID,
			})
		}
	}

	return nil
}
