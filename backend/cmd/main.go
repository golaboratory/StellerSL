package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/user/stellersl/backend/internal/api"
	"github.com/user/stellersl/backend/internal/db"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/stellersl?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	router := chi.NewRouter()
	config := huma.DefaultConfig("StellerSL API", "1.0.0")
	humaAPI := humachi.New(router, config)

	// Middleware for tenant identification from host
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.Split(r.Host, ":")[0]
			tenant, err := queries.GetTenantByDomain(r.Context(), host)
			
			// For local development on localhost, fallback to default tenant if not found
			tenantID := "00000000-0000-0000-0000-000000000001"
			if err == nil {
				tenantID = tenant.ID
			}
			
			// Inject into context (using the same key defined in internal/api)
			// Note: We need a shared key. I'll use a string for simplicity in main.go
			// but internal/api will expect its specific contextKey.
			// Let's align them.
			ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	api.RegisterRoutes(humaAPI, queries)

	fmt.Println("Server starting on :8888")
	http.ListenAndServe(":8888", router)
}
