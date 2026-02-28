package db

import (
	"context"
	"database/sql"
)

type User struct {
	ID           string
	TenantID     string
	Email        string
	PasswordHash string
	Name         string
}

type Tenant struct {
	ID     string
	Name   string
	Domain string
}

type DashboardStats struct {
	TotalTasks     int64
	PendingTasks   int64
	CompletedTasks int64
}

type DailyActivity struct {
	Date      string
	Created   int64
	Completed int64
}

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) GetUserByEmail(ctx context.Context, tenantID, email string) (*User, error) {
	row := q.db.QueryRowContext(ctx, "SELECT id, tenant_id, email, password_hash, name FROM users WHERE tenant_id = $1 AND email = $2", tenantID, email)
	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Name)
	return &u, err
}

func (q *Queries) GetTenantByDomain(ctx context.Context, domain string) (*Tenant, error) {
	row := q.db.QueryRowContext(ctx, "SELECT id, name, domain FROM tenants WHERE domain = $1", domain)
	var t Tenant
	err := row.Scan(&t.ID, &t.Name, &t.Domain)
	return &t, err
}

func (q *Queries) CreateUser(ctx context.Context, tenantID, email, passwordHash, name string) (*User, error) {
	row := q.db.QueryRowContext(ctx, `
		INSERT INTO users (tenant_id, email, password_hash, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tenant_id, email, password_hash, name`,
		tenantID, email, passwordHash, name)
	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Name)
	return &u, err
}

func (q *Queries) GetDashboardStats(ctx context.Context, tenantID, userID string) (*DashboardStats, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total_tasks,
			COUNT(*) FILTER (WHERE status != 'done') as pending_tasks,
			COUNT(*) FILTER (WHERE status = 'done') as completed_tasks
		FROM tasks 
		WHERE tenant_id = $1 AND (assigned_to = $2 OR assigned_to IS NULL)`, tenantID, userID)
	var s DashboardStats
	err := row.Scan(&s.TotalTasks, &s.PendingTasks, &s.CompletedTasks)
	return &s, err
}

func (q *Queries) GetDailyActivity(ctx context.Context, tenantID, userID string) ([]DailyActivity, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT 
			logged_at::text,
			COUNT(*) FILTER (WHERE action = 'task_created') as created_count,
			COUNT(*) FILTER (WHERE action = 'task_completed') as completed_count
		FROM activity_logs
		WHERE tenant_id = $1 AND user_id = $2 AND logged_at > CURRENT_DATE - INTERVAL '7 days'
		GROUP BY logged_at
		ORDER BY logged_at ASC`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []DailyActivity
	for rows.Next() {
		var a DailyActivity
		if err := rows.Scan(&a.Date, &a.Created, &a.Completed); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}
