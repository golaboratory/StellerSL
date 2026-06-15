package task

import (
	"context"
	"database/sql"
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
		engine:  gamification.NewBadgeEngine(),
		streak:  gamification.NewStreakTracker(),
	}
}

func formatNullTime(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format(time.RFC3339)
}

func nullUUIDString(u uuid.NullUUID) string {
	if !u.Valid {
		return ""
	}
	return u.UUID.String()
}

func parseNullTime(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{}
	}
	if parsed, err := time.Parse(time.RFC3339, s); err == nil {
		return sql.NullTime{Time: parsed, Valid: true}
	}
	return sql.NullTime{}
}

func toTaskItem(t db.Task) TaskItem {
	return TaskItem{
		ID:          t.ID.String(),
		ProjectID:   nullUUIDString(t.ProjectID),
		AssignedTo:  nullUUIDString(t.AssignedTo),
		Title:       t.Title,
		Description: t.Description.String,
		Status:      t.Status,
		Priority:    int(t.Priority.Int32),
		DueDate:     formatNullTime(t.DueDate),
		CreatedAt:   formatNullTime(t.CreatedAt),
	}
}

func (s *Service) List(ctx context.Context, tenantID string, input ListTasksInput) (*TaskListOutput, error) {
	params := db.ListTasksParams{
		ProjectID:   db.ToNullUUID(input.ProjectID),
		DueDateFrom: parseNullTime(input.DueDateFrom),
		DueDateTo:   parseNullTime(input.DueDateTo),
		Limit:       input.Limit,
		Offset:      input.Offset,
	}
	if input.Status != "" {
		params.Status = sql.NullString{String: input.Status, Valid: true}
	}
	if input.Priority >= 0 {
		params.Priority = sql.NullInt32{Int32: input.Priority, Valid: true}
	}

	var tasksList []db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		tasksList, err = q.ListTasks(ctx, params)
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskListOutput{}
	for _, t := range tasksList {
		resp.Body.Items = append(resp.Body.Items, toTaskItem(t))
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
	s.fillTaskOutput(resp, t)
	return resp, nil
}

func (s *Service) fillTaskOutput(resp *TaskOutput, t db.Task) {
	resp.Body.ID = t.ID.String()
	resp.Body.ProjectID = nullUUIDString(t.ProjectID)
	resp.Body.AssignedTo = nullUUIDString(t.AssignedTo)
	resp.Body.Title = t.Title
	resp.Body.Description = t.Description.String
	resp.Body.Status = t.Status
	resp.Body.Priority = int(t.Priority.Int32)
	resp.Body.DueDate = formatNullTime(t.DueDate)
	resp.Body.CreatedAt = formatNullTime(t.CreatedAt)
	resp.Body.UpdatedAt = formatNullTime(t.UpdatedAt)
}

func (s *Service) Create(ctx context.Context, tenantID, userID string, input TaskInput) (*TaskOutput, error) {
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
			DueDate:     parseNullTime(input.Body.DueDate),
		})
		if err != nil {
			return err
		}

		uid := db.ParseUUID(userID)
		_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
			TenantID: db.ParseUUID(tenantID),
			UserID:   uid,
			TaskID:   uuid.NullUUID{UUID: t.ID, Valid: true},
			Action:   "task_created",
		})
		_ = s.engine.AwardBadges(ctx, q, uid, "task_created")
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	s.fillTaskOutput(resp, t)
	return resp, nil
}

// applyCompletion records the completion and updates EXP / level / character /
// streak / badges. The activity log is written first so count-based badge
// conditions include the task(s) just completed; the streak is updated before
// badge evaluation so streak badges fire on the day they are reached.
func (s *Service) applyCompletion(ctx context.Context, q *db.Queries, tid, uid uuid.UUID, exp int32, taskIDs ...uuid.UUID) {
	for _, id := range taskIDs {
		_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
			TenantID: tid, UserID: uid,
			TaskID: uuid.NullUUID{UUID: id, Valid: true},
			Action: "task_completed",
		})
	}
	_ = q.AddExp(ctx, db.AddExpParams{UserID: uid, Exp: exp})
	_ = q.UpdateLevel(ctx, uid)
	_ = gamification.UpdateCharacterType(ctx, q, uid)
	_ = s.streak.UpdateStreak(ctx, q, uid)
	_ = s.engine.AwardBadges(ctx, q, uid, "task_completed", taskIDs...)
}

// applyUncompletion reverses EXP and re-evaluates count-based badges after a
// done -> not-done transition.
func (s *Service) applyUncompletion(ctx context.Context, q *db.Queries, tid, uid, taskID uuid.UUID) {
	_ = q.InsertActivityLog(ctx, db.InsertActivityLogParams{
		TenantID: tid, UserID: uid,
		TaskID: uuid.NullUUID{UUID: taskID, Valid: true},
		Action: "task_uncompleted",
	})
	_ = q.SubtractExp(ctx, db.SubtractExpParams{UserID: uid, Exp: 10})
	_ = q.UpdateLevel(ctx, uid)
	_ = gamification.UpdateCharacterType(ctx, q, uid)
	_ = s.engine.ReEvaluateBadges(ctx, q, uid)
}

func (s *Service) Update(ctx context.Context, tenantID, userID, taskID string, input TaskInput) (*TaskOutput, error) {
	var t db.Task
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		// Get current task to detect status transitions
		oldTask, err := q.GetTask(ctx, db.ParseUUID(taskID))
		if err != nil {
			return err
		}

		t, err = q.UpdateTask(ctx, db.UpdateTaskParams{
			ID:          db.ParseUUID(taskID),
			Title:       input.Body.Title,
			Description: sql.NullString{String: input.Body.Description, Valid: input.Body.Description != ""},
			Status:      input.Body.Status,
			Priority:    sql.NullInt32{Int32: int32(input.Body.Priority), Valid: true},
			DueDate:     parseNullTime(input.Body.DueDate),
			AssignedTo:  db.ToNullUUID(input.Body.AssignedTo),
		})
		if err != nil {
			return err
		}

		uid := db.ParseUUID(userID)
		tid := db.ParseUUID(tenantID)

		if t.Status == "done" && oldTask.Status != "done" {
			s.applyCompletion(ctx, q, tid, uid, 10, t.ID)
		} else if oldTask.Status == "done" && t.Status != "done" {
			s.applyUncompletion(ctx, q, tid, uid, t.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	s.fillTaskOutput(resp, t)
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
			s.applyCompletion(ctx, q, tid, uid, 10, t.ID)
		} else if oldTask.Status == "done" && status != "done" {
			s.applyUncompletion(ctx, q, tid, uid, t.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := &TaskOutput{}
	s.fillTaskOutput(resp, t)
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
		_ = s.engine.AwardBadges(ctx, q, uid, "task_created")
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
			s.applyCompletion(ctx, q, tid, uid, int32(10*len(ids)), uuids...)
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
