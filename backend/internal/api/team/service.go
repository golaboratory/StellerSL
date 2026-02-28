package team

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

func (s *Service) List(ctx context.Context, tenantID string) (*TeamListOutput, error) {
	var teams []db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		teams, err = q.ListTeams(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TeamListOutput{}
	for _, t := range teams {
		resp.Body.Items = append(resp.Body.Items, TeamItem{ID: t.ID, Name: t.Name})
	}
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID, name string) (*TeamItem, error) {
	var t *db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.CreateTeam(ctx, tenantID, name)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &TeamItem{ID: t.ID, Name: t.Name}, nil
}

func (s *Service) AddMember(ctx context.Context, tenantID, teamID, userID, role string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.AddTeamMember(ctx, teamID, userID, role)
	})
}
