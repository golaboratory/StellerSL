package project

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterHandlers(api huma.API, service *Service, getTenantID func(context.Context) string) {
	huma.Register(api, huma.Operation{
		OperationID: "list-projects",
		Method:      http.MethodGet,
		Path:        "/projects",
		Summary:     "List Projects",
	}, func(ctx context.Context, input *struct{}) (*ProjectListOutput, error) {
		tenantID := getTenantID(ctx)
		return service.List(ctx, tenantID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-project",
		Method:      http.MethodPost,
		Path:        "/projects",
		Summary:     "Create Project",
	}, func(ctx context.Context, input *ProjectInput) (*ProjectOutput, error) {
		tenantID := getTenantID(ctx)
		return service.Create(ctx, tenantID, *input)
	})
}
