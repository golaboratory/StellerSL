package gamification

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

// StreakTracker maintains the daily completion streak. The querier is passed
// per call so updates run on the caller's tenant-scoped transaction.
type StreakTracker struct {
	now func() time.Time
}

func NewStreakTracker() *StreakTracker {
	return &StreakTracker{now: time.Now}
}

// dateOnly normalizes a timestamp to its local calendar date.
// time.Truncate(24h) operates on absolute (UTC-epoch) days and shifts the
// day boundary in non-UTC timezones, so calendar fields are compared instead.
func dateOnly(t time.Time) time.Time {
	y, m, d := t.Local().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// UpdateStreak updates the daily task completion streak for a user.
// Should be called within a WithTenant transaction when a task is completed.
func (st *StreakTracker) UpdateStreak(ctx context.Context, q db.Querier, userID uuid.UUID) error {
	today := dateOnly(st.now())

	streak, err := q.GetStreak(ctx, db.GetStreakParams{
		UserID:     userID,
		StreakType: "daily_task_completion",
	})

	if err == sql.ErrNoRows {
		// First ever completion
		return q.UpsertStreak(ctx, db.UpsertStreakParams{
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
	if streak.LastDate.Valid && dateOnly(streak.LastDate.Time).Equal(today) {
		return nil
	}

	yesterday := today.AddDate(0, 0, -1)
	newCount := int32(1) // default: reset streak
	if streak.LastDate.Valid && dateOnly(streak.LastDate.Time).Equal(yesterday) {
		newCount = streak.CurrentCount + 1 // consecutive day
	}

	newMax := streak.MaxCount
	if newCount > newMax {
		newMax = newCount
	}

	return q.UpsertStreak(ctx, db.UpsertStreakParams{
		UserID:       userID,
		StreakType:   "daily_task_completion",
		CurrentCount: newCount,
		MaxCount:     newMax,
		LastDate:     sql.NullTime{Time: today, Valid: true},
	})
}
