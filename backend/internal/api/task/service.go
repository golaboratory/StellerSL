package task

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{conn: conn, queries: queries}
}

func (s *Service) List(ctx context.Context, tenantID string) (*TaskListOutput, error) {
	var tasksList []db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		tasksList, err = q.ListTasks(ctx, uuid.NullUUID{})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskListOutput{}
	for _, t := range tasksList {
		resp.Body.Items = append(resp.Body.Items, TaskItem{
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

func (s *Service) Create(ctx context.Context, tenantID string, input TaskInput) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.CreateTask(ctx, db.CreateTaskParams{
			TenantID:    db.ParseUUID(tenantID),
			ProjectID:   db.ToNullUUID(input.Body.ProjectID),
			AssignedTo:  db.ToNullUUID(input.Body.AssignedTo),
			Title:       input.Body.Title,
			Description: sql.NullString{String: input.Body.Description, Valid: input.Body.Description != ""},
			Status:      input.Body.Status,
			Priority:    sql.NullInt32{Int32: int32(input.Body.Priority), Valid: true},
			DueDate:     sql.NullTime{}, // TODO: Parse from input
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	resp.Body.ID = t.ID.String()
	resp.Body.Title = t.Title
	resp.Body.Status = t.Status
	return resp, nil
}

func (s *Service) UpdateStatus(ctx context.Context, tenantID, userID, taskID, status string) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
			ID:     db.ParseUUID(taskID),
			Status: status,
		})
		if err != nil {
			return err
		}

		if status == "done" {
			return q.AddExp(ctx, db.AddExpParams{
				UserID: db.ParseUUID(userID),
				Exp:    10,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	resp.Body.ID = t.ID.String()
	resp.Body.Title = t.Title
	resp.Body.Status = t.Status
	return resp, nil
}

func (s *Service) BulkCreate(ctx context.Context, tenantID string, input BulkTaskCreateInput) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		tid := db.ParseUUID(tenantID)
		for _, t := range input.Body.Tasks {
			_, err := q.CreateTask(ctx, db.CreateTaskParams{
				TenantID:    tid,
				ProjectID:   db.ToNullUUID(t.ProjectID),
				Title:       t.Title,
				Description: sql.NullString{String: t.Description, Valid: t.Description != ""},
				Status:      "todo",
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) BulkUpdateStatus(ctx context.Context, tenantID, userID, status string, ids []string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		uuids := make([]uuid.UUID, len(ids))
		for i, id := range ids {
			uuids[i] = db.ParseUUID(id)
		}

		err := q.BulkUpdateTasksStatus(ctx, db.BulkUpdateTasksStatusParams{
			Column1: uuids,
			Status:  status,
		})
		if err != nil {
			return err
		}

		if status == "done" {
			return q.AddExp(ctx, db.AddExpParams{
				UserID: db.ParseUUID(userID),
				Exp:    int32(10 * len(ids)),
			})
		}
		return nil
	})
}

func (s *Service) BulkDelete(ctx context.Context, tenantID string, ids []string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		uuids := make([]uuid.UUID, len(ids))
		for i, id := range ids {
			uuids[i] = db.ParseUUID(id)
		}
		return q.BulkDeleteTasks(ctx, uuids)
	})
}

func (s *Service) Delete(ctx context.Context, tenantID, taskID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.DeleteTask(ctx, db.ParseUUID(taskID))
	})
}
