package gamification

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

type StreakTracker struct {
	queries *db.Queries
}

func NewStreakTracker(queries *db.Queries) *StreakTracker {
	return &StreakTracker{queries: queries}
}

// UpdateStreak updates the daily task completion streak for a user.
// Should be called within a WithTenant transaction when a task is completed.
func (st *StreakTracker) UpdateStreak(ctx context.Context, userID uuid.UUID) error {
	today := time.Now().Local().Truncate(24 * time.Hour)

	streak, err := st.queries.GetStreak(ctx, db.GetStreakParams{
		UserID:     userID,
		StreakType: "daily_task_completion",
	})

	if err == sql.ErrNoRows {
		// First ever completion
		return st.queries.UpsertStreak(ctx, db.UpsertStreakParams{
			UserID:       userID,
			StreakType:   "daily_task_completion",
			CurrentCount: 1,
			MaxCount:     1,
			LastDate:     sql.NullTime{Time: today, Valid: true},
		})
	}
	if err != nil {
		return err
	}

	// Already counted today
	if streak.LastDate.Valid && streak.LastDate.Time.Truncate(24*time.Hour).Equal(today) {
		return nil
	}

	yesterday := today.AddDate(0, 0, -1)
	newCount := int32(1) // default: reset streak
	if streak.LastDate.Valid && streak.LastDate.Time.Truncate(24*time.Hour).Equal(yesterday) {
		newCount = streak.CurrentCount + 1 // consecutive day
	}

	newMax := streak.MaxCount
	if newCount > newMax {
		newMax = newCount
	}

	return st.queries.UpsertStreak(ctx, db.UpsertStreakParams{
		UserID:       userID,
		StreakType:   "daily_task_completion",
		CurrentCount: newCount,
		MaxCount:     newMax,
		LastDate:     sql.NullTime{Time: today, Valid: true},
	})
}
