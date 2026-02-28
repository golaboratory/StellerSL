package gamification

import (
	"context"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGamificationRouting(t *testing.T) {
	_, api := humatest.New(t)

	RegisterHandlers(api, nil, func(ctx context.Context) (AuthInfo, error) {
		return AuthInfo{TenantID: "tenant-1", UserID: "user-1"}, nil
	})

	spec := api.OpenAPI()

	paths := []string{
		"/growth",
		"/badges",
	}

	for _, p := range paths {
		if spec.Paths[p] == nil {
			t.Errorf("expected path %s to be registered", p)
		}
	}

	if spec.Paths["/growth"].Get == nil || spec.Paths["/growth"].Get.OperationID != "get-user-growth" {
		t.Error("expected get-user-growth operation on GET /growth")
	}

	if spec.Paths["/badges"].Get == nil || spec.Paths["/badges"].Get.OperationID != "list-badges" {
		t.Error("expected list-badges operation on GET /badges")
	}
}
