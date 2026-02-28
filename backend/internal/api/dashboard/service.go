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
	var stats *db.DashboardStats
	var activities []db.DailyActivity

	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		stats, err = q.GetDashboardStats(ctx, userID)
		if err != nil {
			return err
		}

		activities, err = q.GetDailyActivity(ctx, userID)
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
			Created:   a.Created,
			Completed: a.Completed,
		})
	}

	return resp, nil
}
