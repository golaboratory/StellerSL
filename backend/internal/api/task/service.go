package task

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) List(ctx context.Context, tenantID string) (*TaskListOutput, error) {
	tasks, err := s.queries.ListTasks(ctx, tenantID, sql.NullString{})
	if err != nil {
		return nil, err
	}

	resp := &TaskListOutput{}
	for _, t := range tasks {
		resp.Body.Items = append(resp.Body.Items, TaskItem{
			ID:        t.ID,
			ProjectID: t.ProjectID.String,
			Title:     t.Title,
			Status:    t.Status,
			Priority:  t.Priority,
			DueDate:   fmt.Sprintf("%v", t.DueDate.Time),
		})
	}
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID string, input TaskInput) (*TaskOutput, error) {
	var dueDate sql.NullTime
	t, err := s.queries.CreateTask(ctx, tenantID, input.Body.ProjectID, input.Body.AssignedTo, input.Body.Title, input.Body.Description, input.Body.Status, input.Body.Priority, dueDate)
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	resp.Body.ID = t.ID
	resp.Body.Title = t.Title
	resp.Body.Status = t.Status
	return resp, nil
}

func (s *Service) UpdateStatus(ctx context.Context, tenantID, userID, taskID, status string) (*TaskOutput, error) {
	t, err := s.queries.UpdateTaskStatus(ctx, tenantID, taskID, status)
	if err != nil {
		return nil, err
	}

	if status == "done" {
		_ = s.queries.AddExp(ctx, userID, 10)
	}

	resp := &TaskOutput{}
	resp.Body.ID = t.ID
	resp.Body.Title = t.Title
	resp.Body.Status = t.Status
	return resp, nil
}

func (s *Service) BulkCreate(ctx context.Context, tenantID string, input BulkTaskCreateInput) error {
	var dbTasks []db.Task
	for _, t := range input.Body.Tasks {
		dbTasks = append(dbTasks, db.Task{
			ProjectID:   sql.NullString{String: t.ProjectID, Valid: t.ProjectID != ""},
			Title:       t.Title,
			Description: sql.NullString{String: t.Description, Valid: t.Description != ""},
			Status:      "todo",
		})
	}
	return s.queries.BulkCreateTasks(ctx, tenantID, dbTasks)
}

func (s *Service) BulkUpdateStatus(ctx context.Context, tenantID, userID, status string, ids []string) error {
	err := s.queries.BulkUpdateTasksStatus(ctx, tenantID, ids, status)
	if err != nil {
		return err
	}

	if status == "done" {
		_ = s.queries.AddExp(ctx, userID, 10*len(ids))
	}
	return nil
}

func (s *Service) BulkDelete(ctx context.Context, tenantID string, ids []string) error {
	return s.queries.BulkDeleteTasks(ctx, tenantID, ids)
}

func (s *Service) Delete(ctx context.Context, tenantID, taskID string) error {
	return s.queries.DeleteTask(ctx, tenantID, taskID)
}
