package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// WithTenant sets the tenant ID for the current session and returns a new Queries instance
// This is used to enforce Row Level Security (RLS)
func (q *Queries) WithTenant(ctx context.Context, db *sql.DB, tenantID string, fn func(*Queries) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set the session variable for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)); err != nil {
		return err
	}

	if err := fn(New(tx)); err != nil {
		return err
	}

	return tx.Commit()
}

func ParseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}

func ToNullUUID(s string) uuid.NullUUID {
	if s == "" {
		return uuid.NullUUID{Valid: false}
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: id, Valid: true}
}
