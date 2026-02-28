package dashboard

import (
	"context"
	"database/sql"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{conn: conn, queries: queries}
}

func (s *Service) GetStats(ctx context.Context, tenantID, userID string) (*DashboardOutput, error) {
	var stats db.GetDashboardStatsRow
	var activities []db.GetDailyActivityRow
	var recent []db.GetRecentActivityRow

	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		uid := db.ParseUUID(userID)
		stats, err = q.GetDashboardStats(ctx, db.ToNullUUID(userID))
		if err != nil {
			return err
		}

		activities, err = q.GetDailyActivity(ctx, uid)
		if err != nil {
			return err
		}

		recent, err = q.GetRecentActivity(ctx, uid)
		return err
	})

	if err != nil {
		return nil, err
	}

	resp := &DashboardOutput{}
	resp.Body.TotalTasks = stats.TotalTasks
	resp.Body.PendingTasks = stats.PendingTasks
	resp.Body.CompletedTasks = stats.CompletedTasks

	for _, a := range activities {
		resp.Body.DailyActivity = append(resp.Body.DailyActivity, struct {
			Date      string `json:"date"`
			Created   int64  `json:"created"`
			Completed int64  `json:"completed"`
		}{
			Date:      a.Date,
			Created:   a.CreatedCount,
			Completed: a.CompletedCount,
		})
	}

	for _, r := range recent {
		resp.Body.RecentActivity = append(resp.Body.RecentActivity, struct {
			ID        int64  `json:"id"`
			Action    string `json:"action"`
			Date      string `json:"date"`
			TaskTitle string `json:"task_title"`
		}{
			ID:        r.ID,
			Action:    r.Action,
			Date:      r.Date,
			TaskTitle: r.TaskTitle,
		})
	}

	return resp, nil
}
