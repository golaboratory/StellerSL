package project

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/api/task"
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
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description.String,
		})
	}
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID string, input ProjectInput) (*ProjectOutput, error) {
	var p db.Project
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		p, err = q.CreateProject(ctx, db.CreateProjectParams{
			TenantID:    db.ParseUUID(tenantID),
			Name:        input.Body.Name,
			Description: sql.NullString{String: input.Body.Description, Valid: input.Body.Description != ""},
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &ProjectOutput{}
	resp.Body.ID = p.ID.String()
	resp.Body.Name = p.Name
	resp.Body.Description = p.Description.String
	resp.Body.CreatedAt = fmt.Sprintf("%v", p.CreatedAt.Time)
	resp.Body.UpdatedAt = fmt.Sprintf("%v", p.UpdatedAt.Time)
	return resp, nil
}

func (s *Service) Get(ctx context.Context, tenantID, id string) (*ProjectOutput, error) {
	var p db.Project
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		p, err = q.GetProject(ctx, db.ParseUUID(id))
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &ProjectOutput{}
	resp.Body.ID = p.ID.String()
	resp.Body.Name = p.Name
	resp.Body.Description = p.Description.String
	resp.Body.CreatedAt = fmt.Sprintf("%v", p.CreatedAt.Time)
	resp.Body.UpdatedAt = fmt.Sprintf("%v", p.UpdatedAt.Time)
	return resp, nil
}

func (s *Service) Update(ctx context.Context, tenantID, id string, input ProjectInput) (*ProjectOutput, error) {
	var p db.Project
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		p, err = q.UpdateProject(ctx, db.UpdateProjectParams{
			ID:          db.ParseUUID(id),
			Name:        input.Body.Name,
			Description: sql.NullString{String: input.Body.Description, Valid: input.Body.Description != ""},
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &ProjectOutput{}
	resp.Body.ID = p.ID.String()
	resp.Body.Name = p.Name
	resp.Body.Description = p.Description.String
	resp.Body.CreatedAt = fmt.Sprintf("%v", p.CreatedAt.Time)
	resp.Body.UpdatedAt = fmt.Sprintf("%v", p.UpdatedAt.Time)
	return resp, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.DeleteProject(ctx, db.ParseUUID(id))
	})
}

func (s *Service) ListTasks(ctx context.Context, tenantID, id string) (*ProjectTaskListOutput, error) {
	var tasksList []db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		tasksList, err = q.ListTasks(ctx, uuid.NullUUID{UUID: db.ParseUUID(id), Valid: true})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &ProjectTaskListOutput{}
	for _, t := range tasksList {
		resp.Body.Items = append(resp.Body.Items, task.TaskItem{
			ID:        t.ID.String(),
			ProjectID: t.ProjectID.UUID.String(),
			Title:     t.Title,
			Status:    t.Status,
			Priority:  int(t.Priority.Int32),
			DueDate:   fmt.Sprintf("%v", t.DueDate.Time),
		})
	}
	return resp, nil
}

func (s *Service) AssignUser(ctx context.Context, tenantID, projectID, userID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.AssignProjectUser(ctx, db.AssignProjectUserParams{
			ProjectID: db.ParseUUID(projectID),
			UserID:    db.ParseUUID(userID),
		})
	})
}

func (s *Service) UnassignUser(ctx context.Context, tenantID, projectID, userID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.UnassignProjectUser(ctx, db.UnassignProjectUserParams{
			ProjectID: db.ParseUUID(projectID),
			UserID:    db.ParseUUID(userID),
		})
	})
}

func (s *Service) ListMembers(ctx context.Context, tenantID, projectID string) (*task.AccountUserListOutput, error) {
	var users []db.User
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		users, err = q.ListProjectMembers(ctx, db.ParseUUID(projectID))
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &task.AccountUserListOutput{}
	for _, u := range users {
		resp.Body.Items = append(resp.Body.Items, task.AccountUser{
			ID:    u.ID.String(),
			Email: u.Email,
			Name:  u.Name,
		})
	}
	return resp, nil
}
