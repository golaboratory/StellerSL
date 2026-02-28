package project

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

func (s *Service) List(ctx context.Context, tenantID string) (*ProjectListOutput, error) {
	var projects []db.Project
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		projects, err = q.ListProjects(ctx)
		return err
	})
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
	var p *db.Project
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		p, err = q.CreateProject(ctx, tenantID, input.Body.Name, input.Body.Description)
		return err
	})
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
