package db

import (
	"context"
	"database/sql"
)

type Task struct {
	ID          string
	TenantID    string
	ProjectID   sql.NullString
	AssignedTo  sql.NullString
	Title       string
	Description sql.NullString
	Status      string
	Priority    int
	DueDate     sql.NullTime
	CreatedAt   string
	UpdatedAt   string
}

func (q *Queries) CreateTask(ctx context.Context, tenantID, projectID, assignedTo, title, description, status string, priority int, dueDate sql.NullTime) (*Task, error) {
	row := q.db.QueryRowContext(ctx, `
		INSERT INTO tasks (tenant_id, project_id, assigned_to, title, description, status, priority, due_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, tenant_id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at`,
		tenantID, 
		sql.NullString{String: projectID, Valid: projectID != ""},
		sql.NullString{String: assignedTo, Valid: assignedTo != ""},
		title, 
		sql.NullString{String: description, Valid: description != ""},
		status, priority, dueDate)

	var t Task
	err := row.Scan(&t.ID, &t.TenantID, &t.ProjectID, &t.AssignedTo, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func (q *Queries) ListTasks(ctx context.Context, projectID sql.NullString) ([]Task, error) {
	var rows *sql.Rows
	var err error

	if projectID.Valid {
		rows, err = q.db.QueryContext(ctx, `
			SELECT id, tenant_id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at 
			FROM tasks 
			WHERE project_id = $1 
			ORDER BY created_at DESC`, projectID.String)
	} else {
		rows, err = q.db.QueryContext(ctx, `
			SELECT id, tenant_id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at 
			FROM tasks 
			ORDER BY created_at DESC`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.TenantID, &t.ProjectID, &t.AssignedTo, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (q *Queries) GetTask(ctx context.Context, id string) (*Task, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at 
		FROM tasks 
		WHERE id = $1`, id)
	
	var t Task
	err := row.Scan(&t.ID, &t.TenantID, &t.ProjectID, &t.AssignedTo, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func (q *Queries) UpdateTask(ctx context.Context, id, title, description, status string, priority int, dueDate sql.NullTime) (*Task, error) {
	row := q.db.QueryRowContext(ctx, `
		UPDATE tasks 
		SET title = $2, description = $3, status = $4, priority = $5, due_date = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 
		RETURNING id, tenant_id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at`,
		id, title, sql.NullString{String: description, Valid: description != ""}, status, priority, dueDate)

	var t Task
	err := row.Scan(&t.ID, &t.TenantID, &t.ProjectID, &t.AssignedTo, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func (q *Queries) DeleteTask(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, id, status string) (*Task, error) {
	row := q.db.QueryRowContext(ctx, `
		UPDATE tasks 
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 
		RETURNING id, tenant_id, project_id, assigned_to, title, description, status, priority, due_date, created_at, updated_at`,
		id, status)

	var t Task
	err := row.Scan(&t.ID, &t.TenantID, &t.ProjectID, &t.AssignedTo, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func (q *Queries) BulkCreateTasks(ctx context.Context, tenantID string, tasks []Task) error {
	// Assume this is called within WithTenant's transaction
	stmt, err := q.db.PrepareContext(ctx, "INSERT INTO tasks (tenant_id, project_id, assigned_to, title, description, status, priority, due_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range tasks {
		_, err := stmt.ExecContext(ctx, tenantID, t.ProjectID, t.AssignedTo, t.Title, t.Description, t.Status, t.Priority, t.DueDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *Queries) BulkUpdateTasksStatus(ctx context.Context, ids []string, status string) error {
	stmt, err := q.db.PrepareContext(ctx, "UPDATE tasks SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		_, err := stmt.ExecContext(ctx, id, status)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *Queries) BulkDeleteTasks(ctx context.Context, ids []string) error {
	stmt, err := q.db.PrepareContext(ctx, "DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		_, err := stmt.ExecContext(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}
