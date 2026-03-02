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

	// Middleware for tenant identification from host (MUST BE BEFORE humachi.New)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Request: %s %s Host: %s X-Tenant-Host: %s\n", r.Method, r.URL.Path, r.Host, r.Header.Get("X-Tenant-Host"))
			host := strings.Split(r.Host, ":")[0]
			if tenantHost := r.Header.Get("X-Tenant-Host"); tenantHost != "" {
				host = tenantHost
			}
			tenant, err := queries.GetTenantByDomain(r.Context(), host)
			
			// For local development on localhost, fallback to default tenant if not found
			tenantID := "00000000-0000-0000-0000-000000000001"
			if err == nil {
				tenantID = tenant.ID.String()
			}
			
			ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	config := huma.DefaultConfig("StellerSL API", "1.0.0")
	humaAPI := humachi.New(router, config)

	// Serve uploaded files
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		_ = os.Mkdir(uploadDir, 0755)
	}
	fileServer := http.FileServer(http.Dir(uploadDir))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	api.RegisterRoutes(humaAPI, conn)

	fmt.Println("Server starting on :8888")
	if err := http.ListenAndServe(":8888", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
