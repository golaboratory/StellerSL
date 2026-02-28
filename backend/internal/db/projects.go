package db

import (
	"context"
	"database/sql"
)

type Project struct {
	ID          string
	TenantID    string
	Name        string
	Description sql.NullString
	CreatedAt   string
	UpdatedAt   string
}

func (q *Queries) CreateProject(ctx context.Context, tenantID, name, description string) (*Project, error) {
	row := q.db.QueryRowContext(ctx, `
		INSERT INTO projects (tenant_id, name, description) 
		VALUES ($1, $2, $3) 
		RETURNING id, tenant_id, name, description, created_at, updated_at`,
		tenantID, name, sql.NullString{String: description, Valid: description != ""})

	var p Project
	err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (q *Queries) ListProjects(ctx context.Context) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, description, created_at, updated_at 
		FROM projects 
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (q *Queries) GetProject(ctx context.Context, id string) (*Project, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, created_at, updated_at 
		FROM projects 
		WHERE id = $1`, id)
	
	var p Project
	err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (q *Queries) UpdateProject(ctx context.Context, id, name, description string) (*Project, error) {
	row := q.db.QueryRowContext(ctx, `
		UPDATE projects 
		SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 
		RETURNING id, tenant_id, name, description, created_at, updated_at`,
		id, name, sql.NullString{String: description, Valid: description != ""})

	var p Project
	err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (q *Queries) DeleteProject(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}
