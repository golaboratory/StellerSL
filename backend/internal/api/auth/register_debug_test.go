package auth

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/user/stellersl/backend/internal/db"
	_ "github.com/lib/pq"
)

func TestRegisterDebug(t *testing.T) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/stellersl?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	svc := NewService(conn, queries, []byte("secret"))

	input := RegisterInput{}
	input.Body.TenantID = "00000000-0000-0000-0000-000000000001"
	input.Body.Email = "debug-test@example.com"
	input.Body.Password = "password123"
	input.Body.Name = "Debug User"

	err = svc.Register(context.Background(), input)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}
}
