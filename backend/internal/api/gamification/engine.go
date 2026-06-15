package gamification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

// BadgeEngine evaluates badge conditions.
// The querier is passed per call so evaluation runs on the caller's
// tenant-scoped transaction (RLS applies); the engine itself is stateless
// apart from the injectable clock used by time-of-day badges.
type BadgeEngine struct {
	now func() time.Time
}

func NewBadgeEngine() *BadgeEngine {
	return &BadgeEngine{now: time.Now}
}

// AwardBadges evaluates all badge conditions and awards earned ones.
// taskIDs are the tasks that triggered the evaluation (used for
// quick_strike / lightning_fast / zero_inbox checks).
func (e *BadgeEngine) AwardBadges(ctx context.Context, q db.Querier, userID uuid.UUID, trigger string, taskIDs ...uuid.UUID) error {
	// 1. Get current badges to avoid duplicates
	currentBadges, err := q.ListUserBadges(ctx, userID)
	if err != nil {
		return err
	}

	badgeMap := make(map[string]bool)
	for _, b := range currentBadges {
		badgeMap[b.RequirementType] = true
	}

	// 2. Get stats for checking conditions
	stats, err := q.GetDashboardStats(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return err
	}

	// 3. Get all badge definitions
	allBadges, err := q.ListAllBadges(ctx)
	if err != nil {
		return err
	}

	// Tasks that triggered this evaluation (for per-task conditions)
	var triggerTasks []db.Task
	if trigger == "task_completed" {
		for _, id := range taskIDs {
			if t, err := q.GetTask(ctx, id); err == nil {
				triggerTasks = append(triggerTasks, t)
			}
		}
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
				count, err := q.CountTasksCreatedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 5
				}
			}

		case "first_project":
			count, err := q.CountProjectsByUser(ctx, userID)
			if err == nil {
				shouldAward = count >= 1
			}

		case "first_team_join":
			count, err := q.CountTeamMembershipsByUser(ctx, userID)
			if err == nil {
				shouldAward = count >= 1
			}

		// === コンボ・継続バッジ ===
		case "daily_runner_bronze":
			shouldAward = e.checkStreak(ctx, q, userID, 3)

		case "daily_runner_silver":
			shouldAward = e.checkStreak(ctx, q, userID, 7)

		case "daily_runner_gold":
			shouldAward = e.checkStreakMax(ctx, q, userID, 30)

		case "daily_runner_platinum":
			shouldAward = e.checkStreakMax(ctx, q, userID, 100)

		case "comeback_hero":
			if trigger == "task_completed" {
				lastDate, err := q.GetLastActivityDate(ctx, userID)
				// 1970-01-01 is the sentinel for "no previous activity":
				// a first-ever completion is not a comeback.
				if err == nil && lastDate.Year() > 1970 {
					daysSince := int(e.now().Sub(lastDate).Hours() / 24)
					shouldAward = daysSince >= 7
				}
			}

		case "weekend_master":
			count, err := q.CountTasksCompletedOnWeekends(ctx, userID)
			if err == nil {
				shouldAward = count >= 10
			}

		case "monday_motivator":
			if trigger == "task_completed" && e.now().Weekday() == time.Monday {
				count, err := q.CountTasksCompletedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 5
				}
			}

		case "friday_finisher":
			if trigger == "task_completed" && e.now().Weekday() == time.Friday {
				count, err := q.CountTasksCompletedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 5
				}
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
				hour := e.now().Hour()
				shouldAward = hour >= 4 && hour < 9
			}

		case "midnight_owl":
			if trigger == "task_completed" {
				hour := e.now().Hour()
				shouldAward = hour >= 0 && hour < 4
			}

		case "quick_strike":
			for _, t := range triggerTasks {
				if e.now().Sub(t.CreatedAt.Time).Hours() < 1 {
					shouldAward = true
					break
				}
			}

		case "lightning_fast":
			for _, t := range triggerTasks {
				if e.now().Sub(t.CreatedAt.Time).Minutes() < 10 {
					shouldAward = true
					break
				}
			}

		case "all_nighter":
			if trigger == "task_completed" {
				count, err := q.CountTasksCompletedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 10
				}
			}

		case "lunch_hustler":
			if trigger == "task_completed" {
				count, err := q.CountTasksCompletedTodayInHour(ctx, db.CountTasksCompletedTodayInHourParams{
					UserID:  userID,
					Column2: 12,
				})
				if err == nil {
					shouldAward = count >= 3
				}
			}

		case "golden_hour":
			if trigger == "task_completed" {
				count, err := q.CountTasksCompletedTodayInHour(ctx, db.CountTasksCompletedTodayInHourParams{
					UserID:  userID,
					Column2: 17,
				})
				if err == nil {
					shouldAward = count >= 3
				}
			}

		case "hat_trick":
			if trigger == "task_completed" {
				count, err := q.CountTasksCompletedLastHour(ctx, userID)
				if err == nil {
					shouldAward = count >= 3
				}
			}

		// === プロジェクト・整理バッジ ===
		case "organizer":
			count, err := q.CountProjectsByUser(ctx, userID)
			if err == nil {
				shouldAward = count >= 3
			}

		case "multi_tasker":
			if trigger == "task_completed" {
				count, err := q.CountDistinctProjectsCompletedToday(ctx, userID)
				if err == nil {
					shouldAward = count >= 3
				}
			}

		case "zero_inbox":
			for _, t := range triggerTasks {
				if !t.ProjectID.Valid {
					continue
				}
				counts, err := q.GetProjectTaskCounts(ctx, t.ProjectID)
				if err == nil && counts.Total > 0 && counts.Remaining == 0 {
					shouldAward = true
					break
				}
			}
		}

		if shouldAward {
			_ = q.AwardBadge(ctx, db.AwardBadgeParams{
				UserID:  userID,
				BadgeID: b.ID,
			})
		}
	}

	return nil
}

// permanentBadges are never revoked: time-of-day / weekday / single-moment
// achievements remain facts even if the task count later drops
// (バッジの仕様.md: 時間系・ストリーク系・イベント系は永続).
var permanentBadges = map[string]bool{
	"early_bird":            true,
	"midnight_owl":          true,
	"quick_strike":          true,
	"lightning_fast":        true,
	"comeback_hero":         true,
	"all_nighter":           true,
	"lunch_hustler":         true,
	"golden_hour":           true,
	"hat_trick":             true,
	"monday_motivator":      true,
	"friday_finisher":       true,
	"rookie_planner":        true,
	"zero_inbox":            true,
	"first_project":         true,
	"first_team_join":       true,
	"daily_runner_bronze":   true,
	"daily_runner_silver":   true,
	"daily_runner_gold":     true,
	"daily_runner_platinum": true,
}

// ReEvaluateBadges checks if count-based badges should be revoked
// (e.g. after a done -> not-done transition).
func (e *BadgeEngine) ReEvaluateBadges(ctx context.Context, q db.Querier, userID uuid.UUID) error {
	currentBadges, err := q.ListUserBadges(ctx, userID)
	if err != nil {
		return err
	}

	stats, err := q.GetDashboardStats(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return err
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

		if b.RequirementType == "weekend_master" {
			count, err := q.CountTasksCompletedOnWeekends(ctx, userID)
			if err == nil {
				shouldRevoke = count < 10
			}
		}

		if shouldRevoke {
			_ = q.RemoveBadge(ctx, db.RemoveBadgeParams{
				UserID:  userID,
				BadgeID: b.ID,
			})
		}
	}

	return nil
}

func (e *BadgeEngine) checkStreak(ctx context.Context, q db.Querier, userID uuid.UUID, required int32) bool {
	streak, err := q.GetStreak(ctx, db.GetStreakParams{
		UserID:     userID,
		StreakType: "daily_task_completion",
	})
	if err != nil {
		return false
	}
	return streak.CurrentCount >= required
}

func (e *BadgeEngine) checkStreakMax(ctx context.Context, q db.Querier, userID uuid.UUID, required int32) bool {
	streak, err := q.GetStreak(ctx, db.GetStreakParams{
		UserID:     userID,
		StreakType: "daily_task_completion",
	})
	if err != nil {
		return false
	}
	return streak.MaxCount >= required
}
