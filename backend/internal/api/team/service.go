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
		resp.Body.Items = append(resp.Body.Items, TeamItem{ID: t.ID.String(), Name: t.Name})
	}
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID, name string) (*TeamItem, error) {
	var t db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.CreateTeam(ctx, db.CreateTeamParams{
			TenantID: db.ParseUUID(tenantID),
			Name:     name,
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return &TeamItem{ID: t.ID.String(), Name: t.Name}, nil
}

func (s *Service) AddMember(ctx context.Context, tenantID, teamID, userID, role string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.AddTeamMember(ctx, db.AddTeamMemberParams{
			TeamID: db.ParseUUID(teamID),
			UserID: db.ParseUUID(userID),
			Role:   role,
		})
	})
}

func (s *Service) RemoveMember(ctx context.Context, tenantID, teamID, userID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.RemoveTeamMember(ctx, db.RemoveTeamMemberParams{
			TeamID: db.ParseUUID(teamID),
			UserID: db.ParseUUID(userID),
		})
	})
}

func (s *Service) ListMembers(ctx context.Context, tenantID, teamID string) (*TeamMemberListOutput, error) {
	var users []db.User
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		users, err = q.ListTeamMembers(ctx, db.ParseUUID(teamID))
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TeamMemberListOutput{}
	for _, u := range users {
		resp.Body.Items = append(resp.Body.Items, TeamMemberUser{
			ID:    u.ID.String(),
			Email: u.Email,
			Name:  u.Name,
		})
	}
	return resp, nil
}
