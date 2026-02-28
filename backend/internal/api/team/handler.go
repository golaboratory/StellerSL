package team

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	TenantID string
}

func RegisterHandlers(api huma.API, service *Service, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "list-teams",
		Method:      http.MethodGet,
		Path:        "/teams",
		Summary:     "List Teams",
	}, func(ctx context.Context, input *struct{}) (*TeamListOutput, error) {
		auth, _ := getAuth(ctx)
		return service.List(ctx, auth.TenantID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-team",
		Method:      http.MethodPost,
		Path:        "/teams",
		Summary:     "Create Team",
	}, func(ctx context.Context, input *TeamInput) (*TeamItem, error) {
		auth, _ := getAuth(ctx)
		return service.Create(ctx, auth.TenantID, input.Body.Name)
	})

	huma.Register(api, huma.Operation{
		OperationID: "add-team-member",
		Method:      http.MethodPost,
		Path:        "/teams/{id}/members",
		Summary:     "Add Team Member",
	}, func(ctx context.Context, input *TeamMemberInput) (*struct{}, error) {
		if err := service.AddMember(ctx, input.ID, input.Body.UserID, input.Body.Role); err != nil {
			return nil, huma.Error500InternalServerError("Failed to add member")
		}
		return nil, nil
	})
}
