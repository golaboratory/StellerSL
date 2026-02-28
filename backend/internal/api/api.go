package api

import (
	"context"
	"database/sql"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/user/stellersl/backend/internal/api/auth"
	"github.com/user/stellersl/backend/internal/api/dashboard"
	"github.com/user/stellersl/backend/internal/api/gamification"
	"github.com/user/stellersl/backend/internal/api/project"
	"github.com/user/stellersl/backend/internal/api/task"
	"github.com/user/stellersl/backend/internal/api/team"
	"github.com/user/stellersl/backend/internal/db"
)

func getTenantID(ctx context.Context) string {
	if tid, ok := ctx.Value("tenant_id").(string); ok {
		return tid
	}
	return "00000000-0000-0000-0000-000000000001"
}

func getAuth(ctx context.Context) (struct{ TenantID, UserID string }, error) {
	tid := getTenantID(ctx)
	uid, _ := ctx.Value("user_id").(string)
	if uid == "" {
		uid = "00000000-0000-0000-0000-000000000002"
	}
	return struct{ TenantID, UserID string }{TenantID: tid, UserID: uid}, nil
}

func RegisterRoutes(api huma.API, conn *sql.DB) {
	queries := db.New(conn)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("my_secret_key")
	}

	// 1. Initialize Services
	authSvc := auth.NewService(conn, queries, jwtSecret)
	dashSvc := dashboard.NewService(conn, queries)
	projSvc := project.NewService(conn, queries)
	taskSvc := task.NewService(conn, queries)
	teamSvc := team.NewService(conn, queries)
	gamiSvc := gamification.NewService(conn, queries)

	// 2. Register Handlers
	auth.RegisterHandlers(api, authSvc, getTenantID)

	dashboard.RegisterHandlers(api, dashSvc, func(ctx context.Context) (dashboard.AuthInfo, error) {
		a, err := getAuth(ctx)
		return dashboard.AuthInfo{TenantID: a.TenantID, UserID: a.UserID}, err
	})

	project.RegisterHandlers(api, projSvc, getTenantID)

	task.RegisterHandlers(api, taskSvc, func(ctx context.Context) (task.AuthInfo, error) {
		a, err := getAuth(ctx)
		return task.AuthInfo{TenantID: a.TenantID, UserID: a.UserID}, err
	})

	team.RegisterHandlers(api, teamSvc, func(ctx context.Context) (team.AuthInfo, error) {
		a, err := getAuth(ctx)
		return team.AuthInfo{TenantID: a.TenantID}, err
	})

	gamification.RegisterHandlers(api, gamiSvc, func(ctx context.Context) (gamification.AuthInfo, error) {
		a, err := getAuth(ctx)
		return gamification.AuthInfo{UserID: a.UserID}, err
	})
}
