package auth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterHandlers(api huma.API, service *Service, getTenantID func(context.Context) string) {
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/auth/register",
		Summary:     "User Registration",
	}, func(ctx context.Context, input *RegisterInput) (*struct{}, error) {
		if err := service.Register(ctx, *input); err != nil {
			return nil, huma.Error500InternalServerError("Failed to create user")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "User Login",
	}, func(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
		tenantID := getTenantID(ctx)
		resp, err := service.Login(ctx, *input, tenantID)
		if err != nil {
			return nil, huma.Error401Unauthorized("Invalid credentials")
		}
		return resp, nil
	})
}
