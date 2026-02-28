package task

import (
	"context"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestTaskRouting(t *testing.T) {
	_, api := humatest.New(t)

	RegisterHandlers(api, nil, func(ctx context.Context) (AuthInfo, error) {
		return AuthInfo{TenantID: "tenant-1", UserID: "user-1"}, nil
	})

	spec := api.OpenAPI()

	// Verify paths
	paths := []string{
		"/tasks",
		"/tasks/{id}/status",
		"/tasks/bulk",
		"/tasks/bulk/status",
		"/tasks/{id}",
	}

	for _, p := range paths {
		if spec.Paths[p] == nil {
			t.Errorf("expected path %s to be registered", p)
		}
	}

	// Verify specific operations
	if spec.Paths["/tasks"].Get == nil || spec.Paths["/tasks"].Get.OperationID != "list-tasks" {
		t.Error("expected list-tasks operation on GET /tasks")
	}

	if spec.Paths["/tasks"].Post == nil || spec.Paths["/tasks"].Post.OperationID != "create-task" {
		t.Error("expected create-task operation on POST /tasks")
	}
}
