package dashboard

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	TenantID string
	UserID   string
}

func RegisterHandlers(api huma.API, service *Service, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "get-dashboard-stats",
		Method:      http.MethodGet,
		Path:        "/dashboard",
		Summary:     "Dashboard Stats",
	}, func(ctx context.Context, input *struct{}) (*DashboardOutput, error) {
		auth, err := getAuth(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}

		resp, err := service.GetStats(ctx, auth.TenantID, auth.UserID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to fetch stats")
		}
		return resp, nil
	})
}
