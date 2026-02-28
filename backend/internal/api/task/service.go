package task

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/api/gamification"
	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
	engine  *gamification.BadgeEngine
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{
		conn:    conn,
		queries: queries,
		engine:  gamification.NewBadgeEngine(queries),
	}
}

func (s *Service) List(ctx context.Context, tenantID string, limit, offset int32) (*TaskListOutput, error) {
	var tasksList []db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		tasksList, err = q.ListTasks(ctx, db.ListTasksParams{
			ProjectID: uuid.NullUUID{},
			Limit:     limit,
			Offset:    offset,
		})
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
		dueDate := sql.NullTime{}
		if input.Body.DueDate != "" {
			if parsed, err := time.Parse(time.RFC3339, input.Body.DueDate); err == nil {
				dueDate = sql.NullTime{Time: parsed, Valid: true}
			}
		}

		t, err = q.CreateTask(ctx, db.CreateTaskParams{
			TenantID:    db.ParseUUID(tenantID),
			ProjectID:   db.ToNullUUID(input.Body.ProjectID),
			AssignedTo:  db.ToNullUUID(input.Body.AssignedTo),
			Title:       input.Body.Title,
			Description: sql.NullString{String: input.Body.Description, Valid: input.Body.Description != ""},
			Status:      input.Body.Status,
			Priority:    sql.NullInt32{Int32: int32(input.Body.Priority), Valid: true},
			DueDate:     dueDate,
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

func (s *Service) Update(ctx context.Context, tenantID, userID, taskID string, input TaskInput) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		dueDate := sql.NullTime{}
		if input.Body.DueDate != "" {
			if parsed, err := time.Parse(time.RFC3339, input.Body.DueDate); err == nil {
				dueDate = sql.NullTime{Time: parsed, Valid: true}
			}
		}

		t, err = q.UpdateTask(ctx, db.UpdateTaskParams{
			ID:          db.ParseUUID(taskID),
			Title:       input.Body.Title,
			Description: sql.NullString{String: input.Body.Description, Valid: input.Body.Description != ""},
			Status:      input.Body.Status,
			Priority:    sql.NullInt32{Int32: int32(input.Body.Priority), Valid: true},
			DueDate:     dueDate,
		})
		if err != nil {
			return err
		}

		if t.Status == "done" {
			_ = q.AddExp(ctx, db.AddExpParams{UserID: db.ParseUUID(userID), Exp: 10})
			_ = q.UpdateLevel(ctx, db.ParseUUID(userID))
			_ = s.engine.AwardBadges(ctx, db.ParseUUID(userID), "task_completed")
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
			_ = q.AddExp(ctx, db.AddExpParams{UserID: db.ParseUUID(userID), Exp: 10})
			_ = q.UpdateLevel(ctx, db.ParseUUID(userID))
			_ = s.engine.AwardBadges(ctx, db.ParseUUID(userID), "task_completed")
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
			uid := db.ParseUUID(userID)
			_ = q.AddExp(ctx, db.AddExpParams{UserID: uid, Exp: int32(10 * len(ids))})
			_ = q.UpdateLevel(ctx, uid)
			_ = s.engine.AwardBadges(ctx, uid, "task_completed")
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
