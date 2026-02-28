package project

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

func (s *Service) List(ctx context.Context, tenantID string) (*ProjectListOutput, error) {
	projects, err := s.queries.ListProjects(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	resp := &ProjectListOutput{}
	for _, p := range projects {
		resp.Body.Items = append(resp.Body.Items, ProjectItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description.String,
		})
	}
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID string, input ProjectInput) (*ProjectOutput, error) {
	p, err := s.queries.CreateProject(ctx, tenantID, input.Body.Name, input.Body.Description)
	if err != nil {
		return nil, err
	}

	resp := &ProjectOutput{}
	resp.Body.ID = p.ID
	resp.Body.Name = p.Name
	resp.Body.Description = p.Description.String
	resp.Body.CreatedAt = p.CreatedAt
	resp.Body.UpdatedAt = p.UpdatedAt
	return resp, nil
}
