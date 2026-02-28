package project

import (
	"context"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestProjectRouting(t *testing.T) {
	_, api := humatest.New(t)

	RegisterHandlers(api, nil, func(ctx context.Context) string {
		return "tenant-1"
	})

	spec := api.OpenAPI()

	paths := []string{
		"/projects",
		"/projects/{id}",
		"/projects/{id}/tasks",
		"/projects/{id}/users",
		"/projects/{id}/users/{user_id}",
	}

	for _, p := range paths {
		if spec.Paths[p] == nil {
			t.Errorf("expected path %s to be registered", p)
		}
	}

	// Verify specific operations
	if spec.Paths["/projects"].Get == nil || spec.Paths["/projects"].Get.OperationID != "list-projects" {
		t.Error("expected list-projects operation on GET /projects")
	}

	if spec.Paths["/projects"].Post == nil || spec.Paths["/projects"].Post.OperationID != "create-project" {
		t.Error("expected create-project operation on POST /projects")
	}

	if spec.Paths["/projects/{id}"].Get == nil || spec.Paths["/projects/{id}"].Get.OperationID != "get-project" {
		t.Error("expected get-project operation on GET /projects/{id}")
	}
}
