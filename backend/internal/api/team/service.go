package team

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

func (s *Service) List(ctx context.Context, tenantID string) (*TeamListOutput, error) {
	teams, err := s.queries.ListTeams(ctx, tenantID)
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
	t, err := s.queries.CreateTeam(ctx, tenantID, name)
	if err != nil {
		return nil, err
	}
	return &TeamItem{ID: t.ID, Name: t.Name}, nil
}

func (s *Service) AddMember(ctx context.Context, teamID, userID, role string) error {
	return s.queries.AddTeamMember(ctx, teamID, userID, role)
}
