package api

import (
	"context"
	"testing"
)

func TestGetTenantID(t *testing.T) {
	ctx := context.WithValue(context.Background(), "tenant_id", "test-tenant")
	tid := getTenantID(ctx)
	if tid != "test-tenant" {
		t.Errorf("expected test-tenant, got %s", tid)
	}

	tidFallback := getTenantID(context.Background())
	if tidFallback != "00000000-0000-0000-0000-000000000001" {
		t.Errorf("expected default fallback, got %s", tidFallback)
	}
}

func TestGetAuth(t *testing.T) {
	ctx := context.WithValue(context.Background(), "tenant_id", "tid123")
	ctx = context.WithValue(ctx, "user_id", "uid456")
	
	auth, err := getAuth(ctx)
	if err != nil {
		t.Fatal(err)
	}
	
	if auth.TenantID != "tid123" {
		t.Errorf("expected tid123, got %s", auth.TenantID)
	}
	if auth.UserID != "uid456" {
		t.Errorf("expected uid456, got %s", auth.UserID)
	}
}
