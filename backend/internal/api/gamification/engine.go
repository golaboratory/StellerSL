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

// AwardBadges evaluates all badge conditions and awards earned ones.
// taskID is optional, used for quick_strike badge (time since task creation).
func (e *BadgeEngine) AwardBadges(ctx context.Context, userID uuid.UUID, trigger string, taskID ...uuid.UUID) error {
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

	// 3. Get all badge definitions
	allBadges, err := e.queries.ListAllBadges(ctx)
	if err != nil {
		return err
	}

	for _, b := range allBadges {
		if badgeMap[b.RequirementType] {
			continue // Already have it
		}

		shouldAward := false
		switch b.RequirementType {

		// === はじまりのバッジ ===
		case "first_task":
			shouldAward = stats.CompletedTasks >= 1

		case "three_strike":
			shouldAward = stats.CompletedTasks >= 3

		case "rookie_planner":
			if trigger == "task_created" {
				count, err := e.queries.CountTasksCreatedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 5
				}
			}

		// === コンボ・継続バッジ ===
		case "daily_runner_bronze":
			shouldAward = e.checkStreak(ctx, userID, 3)

		case "daily_runner_silver":
			shouldAward = e.checkStreak(ctx, userID, 7)

		case "daily_runner_gold":
			shouldAward = e.checkStreakMax(ctx, userID, 30)

		case "comeback_hero":
			if trigger == "task_completed" {
				lastDate, err := e.queries.GetLastActivityDate(ctx, userID)
				if err == nil {
					daysSince := int(time.Since(lastDate).Hours() / 24)
					shouldAward = daysSince >= 7
				}
			}

		case "weekend_master":
			count, err := e.queries.CountTasksCompletedOnWeekends(ctx, userID)
			if err == nil {
				shouldAward = count >= 10
			}

		// === マイルストーンバッジ ===
		case "task_crusher_1":
			shouldAward = stats.CompletedTasks >= 10
		case "task_crusher_2":
			shouldAward = stats.CompletedTasks >= 50
		case "task_crusher_3":
			shouldAward = stats.CompletedTasks >= 100
		case "task_crusher_4":
			shouldAward = stats.CompletedTasks >= 500
		case "task_crusher_5":
			shouldAward = stats.CompletedTasks >= 1000

		// === スピード＆タイムバッジ ===
		case "early_bird":
			if trigger == "task_completed" {
				hour := time.Now().Hour()
				shouldAward = hour >= 4 && hour < 9
			}

		case "midnight_owl":
			if trigger == "task_completed" {
				hour := time.Now().Hour()
				shouldAward = hour >= 0 && hour < 4
			}

		case "quick_strike":
			if trigger == "task_completed" && len(taskID) > 0 {
				task, err := e.queries.GetTask(ctx, taskID[0])
				if err == nil {
					shouldAward = time.Since(task.CreatedAt.Time).Hours() < 1
				}
			}

		case "all_nighter":
			if trigger == "task_completed" {
				count, err := e.queries.CountTasksCompletedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 10
				}
			}

		// === プロジェクト・整理バッジ ===
		case "organizer":
			count, err := e.queries.CountProjectsByUser(ctx, userID)
			if err == nil {
				shouldAward = count >= 3
			}

		case "multi_tasker":
			if trigger == "task_completed" {
				count, err := e.queries.CountDistinctProjectsCompletedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 3
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

// ReEvaluateBadges checks if count-based badges should be revoked.
// Time-based badges (early_bird, midnight_owl) are permanent once earned.
func (e *BadgeEngine) ReEvaluateBadges(ctx context.Context, userID uuid.UUID) error {
	currentBadges, err := e.queries.ListUserBadges(ctx, userID)
	if err != nil {
		return err
	}

	stats, err := e.queries.GetDashboardStats(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return err
	}

	permanentBadges := map[string]bool{
		"early_bird":   true,
		"midnight_owl": true,
		"quick_strike": true,
		"comeback_hero": true,
	}

	thresholds := map[string]int64{
		"first_task":     1,
		"three_strike":   3,
		"task_crusher_1": 10,
		"task_crusher_2": 50,
		"task_crusher_3": 100,
		"task_crusher_4": 500,
		"task_crusher_5": 1000,
	}

	for _, b := range currentBadges {
		if permanentBadges[b.RequirementType] {
			continue
		}

		shouldRevoke := false
		if threshold, ok := thresholds[b.RequirementType]; ok {
			shouldRevoke = stats.CompletedTasks < threshold
		}

		// weekend_master and all_nighter also need count checks
		switch b.RequirementType {
		case "weekend_master":
			count, err := e.queries.CountTasksCompletedOnWeekends(ctx, userID)
			if err == nil {
				shouldRevoke = count < 10
			}
		case "all_nighter":
			count, err := e.queries.CountTasksCompletedToday(ctx, userID)
			if err == nil {
				shouldRevoke = count < 10
			}
		}

		if shouldRevoke {
			_ = e.queries.RemoveBadge(ctx, db.RemoveBadgeParams{
				UserID:  userID,
				BadgeID: b.ID,
			})
		}
	}

	return nil
}

func (e *BadgeEngine) checkStreak(ctx context.Context, userID uuid.UUID, required int32) bool {
	streak, err := e.queries.GetStreak(ctx, db.GetStreakParams{
		UserID:     userID,
		StreakType: "daily_task_completion",
	})
	if err != nil {
		return false
	}
	return streak.CurrentCount >= required
}

func (e *BadgeEngine) checkStreakMax(ctx context.Context, userID uuid.UUID, required int32) bool {
	streak, err := e.queries.GetStreak(ctx, db.GetStreakParams{
		UserID:     userID,
		StreakType: "daily_task_completion",
	})
	if err != nil {
		return false
	}
	return streak.MaxCount >= required
}
