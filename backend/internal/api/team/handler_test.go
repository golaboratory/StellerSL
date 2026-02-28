package team

import (
	"context"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestTeamRouting(t *testing.T) {
	_, api := humatest.New(t)

	RegisterHandlers(api, nil, func(ctx context.Context) (AuthInfo, error) {
		return AuthInfo{TenantID: "tenant-1", UserID: "user-1"}, nil
	})

	spec := api.OpenAPI()

	paths := []string{
		"/teams",
		"/teams/{id}/members",
		"/teams/{id}/members/{user_id}",
	}

	for _, p := range paths {
		if spec.Paths[p] == nil {
			t.Errorf("expected path %s to be registered", p)
		}
	}

	if spec.Paths["/teams"].Get == nil || spec.Paths["/teams"].Get.OperationID != "list-teams" {
		t.Error("expected list-teams operation on GET /teams")
	}

	if spec.Paths["/teams"].Post == nil || spec.Paths["/teams"].Post.OperationID != "create-team" {
		t.Error("expected create-team operation on POST /teams")
	}

	if spec.Paths["/teams/{id}/members"].Get == nil || spec.Paths["/teams/{id}/members"].Get.OperationID != "list-team-members" {
		t.Error("expected list-team-members operation on GET /teams/{id}/members")
	}

	if spec.Paths["/teams/{id}/members"].Post == nil || spec.Paths["/teams/{id}/members"].Post.OperationID != "add-team-member" {
		t.Error("expected add-team-member operation on POST /teams/{id}/members")
	}

	if spec.Paths["/teams/{id}/members/{user_id}"].Delete == nil || spec.Paths["/teams/{id}/members/{user_id}"].Delete.OperationID != "remove-team-member" {
		t.Error("expected remove-team-member operation on DELETE /teams/{id}/members/{user_id}")
	}
}
