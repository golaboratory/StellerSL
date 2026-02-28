package auth

import (
	"context"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestAuthRouting(t *testing.T) {
	_, api := humatest.New(t)

	RegisterHandlers(api, nil, func(ctx context.Context) string {
		return "tenant-1"
	}, func(ctx context.Context) (AuthInfo, error) {
		return AuthInfo{TenantID: "tenant-1", UserID: "user-1"}, nil
	})

	spec := api.OpenAPI()

	paths := []string{
		"/auth/register",
		"/auth/login",
		"/auth/profile",
	}

	for _, p := range paths {
		if spec.Paths[p] == nil {
			t.Errorf("expected path %s to be registered", p)
		}
	}

	if spec.Paths["/auth/register"].Post == nil || spec.Paths["/auth/register"].Post.OperationID != "register" {
		t.Error("expected register operation on POST /auth/register")
	}

	if spec.Paths["/auth/login"].Post == nil || spec.Paths["/auth/login"].Post.OperationID != "login" {
		t.Error("expected login operation on POST /auth/login")
	}

	if spec.Paths["/auth/profile"].Put == nil || spec.Paths["/auth/profile"].Put.OperationID != "update-profile" {
		t.Error("expected update-profile operation on PUT /auth/profile")
	}
}
