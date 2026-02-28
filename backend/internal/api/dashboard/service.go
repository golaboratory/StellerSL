package dashboard

import (
	"context"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) GetStats(ctx context.Context, tenantID, userID string) (*DashboardOutput, error) {
	stats, err := s.queries.GetDashboardStats(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	activities, err := s.queries.GetDailyActivity(ctx, tenantID, userID)
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
