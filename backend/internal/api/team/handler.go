package team

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	TenantID string
	UserID   string
}

type ListTeamsInput struct {
	Limit  int32 `query:"limit" default:"50" maximum:"200"`
	Offset int32 `query:"offset" default:"0"`
}

func RegisterHandlers(api huma.API, service *Service, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "list-teams",
		Method:      http.MethodGet,
		Path:        "/teams",
		Summary:     "List Teams",
	}, func(ctx context.Context, input *ListTeamsInput) (*TeamListOutput, error) {
		auth, _ := getAuth(ctx)
		return service.List(ctx, auth.TenantID, input.Limit, input.Offset)
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
		auth, _ := getAuth(ctx)
		if err := service.AddMember(ctx, auth.TenantID, input.ID, input.Body.UserID, input.Body.Role); err != nil {
			return nil, huma.Error500InternalServerError("Failed to add member")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-team-members",
		Method:      http.MethodGet,
		Path:        "/teams/{id}/members",
		Summary:     "List Team Members",
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
	}) (*TeamMemberListOutput, error) {
		auth, _ := getAuth(ctx)
		return service.ListMembers(ctx, auth.TenantID, input.ID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "remove-team-member",
		Method:      http.MethodDelete,
		Path:        "/teams/{id}/members/{user_id}",
		Summary:     "Remove Team Member",
	}, func(ctx context.Context, input *struct {
		ID     string `path:"id"`
		UserID string `path:"user_id"`
	}) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.RemoveMember(ctx, auth.TenantID, input.ID, input.UserID); err != nil {
			return nil, err
		}
		return nil, nil
	})
}
