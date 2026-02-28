package dashboard

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

// mockService is a stub that doesn't hit the DB for routing tests
func TestDashboardRouting(t *testing.T) {
	_, api := humatest.New(t)

	// We only need to verify that RegisterHandlers doesn't panic
	// and correctly registers the route.
	RegisterHandlers(api, nil, func(ctx context.Context) (AuthInfo, error) {
		return AuthInfo{TenantID: "tenant-1", UserID: "user-1"}, nil
	})

	// Verify operation exists in OpenAPI spec
	spec := api.OpenAPI()
	
	if spec.Paths["/dashboard"] == nil {
		t.Fatal("expected /dashboard path to be registered")
	}
	
	if spec.Paths["/dashboard"].Get == nil {
		t.Fatal("expected GET /dashboard to be registered")
	}

	if spec.Paths["/dashboard"].Get.OperationID != "get-dashboard-stats" {
		t.Errorf("expected operation ID get-dashboard-stats, got %s", spec.Paths["/dashboard"].Get.OperationID)
	}
}

func TestAuthFailure(t *testing.T) {
	// A simple test to verify standard huma error formatting 
	// (usually handled by huma, but good to have a placeholder)
	err := huma.Error401Unauthorized("Unauthorized")
	if err.GetStatus() != http.StatusUnauthorized {
		t.Errorf("expected 401 status, got %d", err.GetStatus())
	}
}
