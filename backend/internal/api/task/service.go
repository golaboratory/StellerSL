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
	streak  *gamification.StreakTracker
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{
		conn:    conn,
		queries: queries,
		engine:  gamification.NewBadgeEngine(queries),
		streak:  gamification.NewStreakTracker(queries),
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

func (s *Service) Get(ctx context.Context, tenantID, taskID string) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.GetTask(ctx, db.ParseUUID(taskID))
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	resp.Body.ID = t.ID.String()
	resp.Body.ProjectID = t.ProjectID.UUID.String()
	resp.Body.Title = t.Title
	resp.Body.Description = t.Description.String
	resp.Body.Status = t.Status
	resp.Body.Priority = int(t.Priority.Int32)
	resp.Body.DueDate = fmt.Sprintf("%v", t.DueDate.Time)
	resp.Body.CreatedAt = t.CreatedAt.Time.Format(time.RFC3339)
	resp.Body.UpdatedAt = t.UpdatedAt.Time.Format(time.RFC3339)
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID, userID string, input TaskInput) (*TaskOutput, error) {
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
		if err != nil {
			return err
		}

		_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
			TenantID: db.ParseUUID(tenantID),
			UserID:   db.ParseUUID(userID),
			TaskID:   uuid.NullUUID{UUID: t.ID, Valid: true},
			Action:   "task_created",
		})
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

func (s *Service) Update(ctx context.Context, tenantID, userID, taskID string, input TaskInput) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		// Get current task to detect status transitions
		oldTask, err := q.GetTask(ctx, db.ParseUUID(taskID))
		if err != nil {
			return err
		}

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

		uid := db.ParseUUID(userID)
		tid := db.ParseUUID(tenantID)

		if t.Status == "done" && oldTask.Status != "done" {
			// Transition to done: award EXP and badges
			_ = q.AddExp(ctx, db.AddExpParams{UserID: uid, Exp: 10})
			_ = q.UpdateLevel(ctx, uid)
			_ = gamification.UpdateCharacterType(ctx, q, uid)
			_ = s.engine.AwardBadges(ctx, uid, "task_completed", t.ID)
			_ = s.streak.UpdateStreak(ctx, uid)
			_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
				TenantID: tid, UserID: uid,
				TaskID: uuid.NullUUID{UUID: t.ID, Valid: true},
				Action: "task_completed",
			})
		} else if oldTask.Status == "done" && t.Status != "done" {
			// Transition from done: reverse EXP and re-evaluate badges
			_ = q.SubtractExp(ctx, db.SubtractExpParams{UserID: uid, Exp: 10})
			_ = q.UpdateLevel(ctx, uid)
			_ = gamification.UpdateCharacterType(ctx, q, uid)
			_ = s.engine.ReEvaluateBadges(ctx, uid)
			_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
				TenantID: tid, UserID: uid,
				TaskID: uuid.NullUUID{UUID: t.ID, Valid: true},
				Action: "task_uncompleted",
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

func (s *Service) UpdateStatus(ctx context.Context, tenantID, userID, taskID, status string) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		// Get current task to detect status transitions
		oldTask, err := q.GetTask(ctx, db.ParseUUID(taskID))
		if err != nil {
			return err
		}

		t, err = q.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
			ID:     db.ParseUUID(taskID),
			Status: status,
		})
		if err != nil {
			return err
		}

		uid := db.ParseUUID(userID)
		tid := db.ParseUUID(tenantID)

		if status == "done" && oldTask.Status != "done" {
			_ = q.AddExp(ctx, db.AddExpParams{UserID: uid, Exp: 10})
			_ = q.UpdateLevel(ctx, uid)
			_ = gamification.UpdateCharacterType(ctx, q, uid)
			_ = s.engine.AwardBadges(ctx, uid, "task_completed", t.ID)
			_ = s.streak.UpdateStreak(ctx, uid)
			_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
				TenantID: tid, UserID: uid,
				TaskID: uuid.NullUUID{UUID: t.ID, Valid: true},
				Action: "task_completed",
			})
		} else if oldTask.Status == "done" && status != "done" {
			_ = q.SubtractExp(ctx, db.SubtractExpParams{UserID: uid, Exp: 10})
			_ = q.UpdateLevel(ctx, uid)
			_ = gamification.UpdateCharacterType(ctx, q, uid)
			_ = s.engine.ReEvaluateBadges(ctx, uid)
			_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
				TenantID: tid, UserID: uid,
				TaskID: uuid.NullUUID{UUID: t.ID, Valid: true},
				Action: "task_uncompleted",
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

func (s *Service) BulkCreate(ctx context.Context, tenantID, userID string, input BulkTaskCreateInput) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		tid := db.ParseUUID(tenantID)
		uid := db.ParseUUID(userID)
		for _, t := range input.Body.Tasks {
			created, err := q.CreateTask(ctx, db.CreateTaskParams{
				TenantID:    tid,
				ProjectID:   db.ToNullUUID(t.ProjectID),
				Title:       t.Title,
				Description: sql.NullString{String: t.Description, Valid: t.Description != ""},
				Status:      "todo",
			})
			if err != nil {
				return err
			}
			_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
				TenantID: tid,
				UserID:   uid,
				TaskID:   uuid.NullUUID{UUID: created.ID, Valid: true},
				Action:   "task_created",
			})
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
			tid := db.ParseUUID(tenantID)
			_ = q.AddExp(ctx, db.AddExpParams{UserID: uid, Exp: int32(10 * len(ids))})
			_ = q.UpdateLevel(ctx, uid)
			_ = gamification.UpdateCharacterType(ctx, q, uid)
			_ = s.engine.AwardBadges(ctx, uid, "task_completed")
			_ = s.streak.UpdateStreak(ctx, uid)
			for _, id := range uuids {
				_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
					TenantID: tid,
					UserID:   uid,
					TaskID:   uuid.NullUUID{UUID: id, Valid: true},
					Action:   "task_completed",
				})
			}
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
