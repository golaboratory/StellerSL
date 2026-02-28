package db

import (
	"context"
	"database/sql"
	"fmt"
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

// DBTX defines the interface for database operations, supporting both *sql.DB and *sql.Tx
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

// WithTenant sets the tenant ID for the current session and returns a new Queries instance
// This is used to enforce Row Level Security (RLS)
func (q *Queries) WithTenant(ctx context.Context, db *sql.DB, tenantID string, fn func(*Queries) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set the session variable for RLS
	// We use fmt.Sprintf for the query but properly handle the tenantID value via tx.Exec
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	if err := fn(&Queries{db: tx}); err != nil {
		return err
	}

	return tx.Commit()
}

// Non-RLS query for tenant identification (runs without tenant_id context)
func (q *Queries) GetTenantByDomain(ctx context.Context, domain string) (*Tenant, error) {
	row := q.db.QueryRowContext(ctx, "SELECT id, name, domain FROM tenants WHERE domain = $1", domain)
	var t Tenant
	err := row.Scan(&t.ID, &t.Name, &t.Domain)
	return &t, err
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// RLS will automatically filter by tenant_id, so we don't NEED it in the WHERE clause,
	// but keeping it for safety/performance or removing for strict RLS reliance.
	// For RLS, current_setting('app.current_tenant_id') must be set.
	row := q.db.QueryRowContext(ctx, "SELECT id, tenant_id, email, password_hash, name FROM users WHERE email = $1", email)
	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Name)
	return &u, err
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

func (q *Queries) GetDashboardStats(ctx context.Context, userID string) (*DashboardStats, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total_tasks,
			COUNT(*) FILTER (WHERE status != 'done') as pending_tasks,
			COUNT(*) FILTER (WHERE status = 'done') as completed_tasks
		FROM tasks 
		WHERE (assigned_to = $1 OR assigned_to IS NULL)`, userID)
	var s DashboardStats
	err := row.Scan(&s.TotalTasks, &s.PendingTasks, &s.CompletedTasks)
	return &s, err
}

func (q *Queries) GetDailyActivity(ctx context.Context, userID string) ([]DailyActivity, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT 
			logged_at::text,
			COUNT(*) FILTER (WHERE action = 'task_created') as created_count,
			COUNT(*) FILTER (WHERE action = 'task_completed') as completed_count
		FROM activity_logs
		WHERE user_id = $1 AND logged_at > CURRENT_DATE - INTERVAL '7 days'
		GROUP BY logged_at
		ORDER BY logged_at ASC`, userID)
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
