package project

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/user/stellersl/backend/internal/api/task"
)

type ListProjectsInput struct {
	Limit  int32 `query:"limit" default:"50" maximum:"200"`
	Offset int32 `query:"offset" default:"0"`
}

type UserSearchInput struct {
	Query string `query:"q" minLength:"1"`
}

func RegisterHandlers(api huma.API, service *Service, getTenantID func(context.Context) string) {
	huma.Register(api, huma.Operation{
		OperationID: "list-projects",
		Method:      http.MethodGet,
		Path:        "/projects",
		Summary:     "List Projects",
	}, func(ctx context.Context, input *ListProjectsInput) (*ProjectListOutput, error) {
		tenantID := getTenantID(ctx)
		return service.List(ctx, tenantID, input.Limit, input.Offset)
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

	huma.Register(api, huma.Operation{
		OperationID: "get-project",
		Method:      http.MethodGet,
		Path:        "/projects/{id}",
		Summary:     "Get Project",
	}, func(ctx context.Context, input *ProjectIDInput) (*ProjectOutput, error) {
		tenantID := getTenantID(ctx)
		return service.Get(ctx, tenantID, input.ID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-project",
		Method:      http.MethodPut,
		Path:        "/projects/{id}",
		Summary:     "Update Project",
	}, func(ctx context.Context, input *struct {
		ProjectIDInput
		Body ProjectInput `json:"body"`
	}) (*ProjectOutput, error) {
		tenantID := getTenantID(ctx)
		return service.Update(ctx, tenantID, input.ID, input.Body)
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-project",
		Method:      http.MethodDelete,
		Path:        "/projects/{id}",
		Summary:     "Delete Project",
	}, func(ctx context.Context, input *ProjectIDInput) (*struct{}, error) {
		tenantID := getTenantID(ctx)
		if err := service.Delete(ctx, tenantID, input.ID); err != nil {
			return nil, err
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-project-tasks",
		Method:      http.MethodGet,
		Path:        "/projects/{id}/tasks",
		Summary:     "List Project Tasks",
	}, func(ctx context.Context, input *ProjectIDInput) (*ProjectTaskListOutput, error) {
		tenantID := getTenantID(ctx)
		return service.ListTasks(ctx, tenantID, input.ID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "assign-project-user",
		Method:      http.MethodPost,
		Path:        "/projects/{id}/users",
		Summary:     "Assign User to Project",
	}, func(ctx context.Context, input *ProjectUserAssignmentInput) (*struct{}, error) {
		tenantID := getTenantID(ctx)
		if err := service.AssignUser(ctx, tenantID, input.ID, input.Body.UserID); err != nil {
			return nil, err
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "unassign-project-user",
		Method:      http.MethodDelete,
		Path:        "/projects/{id}/users/{user_id}",
		Summary:     "Unassign User from Project",
	}, func(ctx context.Context, input *ProjectUserUnassignmentInput) (*struct{}, error) {
		tenantID := getTenantID(ctx)
		if err := service.UnassignUser(ctx, tenantID, input.ID, input.UserID); err != nil {
			return nil, err
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-project-users",
		Method:      http.MethodGet,
		Path:        "/projects/{id}/users",
		Summary:     "List Project Users",
	}, func(ctx context.Context, input *ProjectIDInput) (*task.AccountUserListOutput, error) {
		tenantID := getTenantID(ctx)
		return service.ListMembers(ctx, tenantID, input.ID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "search-users",
		Method:      http.MethodGet,
		Path:        "/users/search",
		Summary:     "Search Users",
	}, func(ctx context.Context, input *UserSearchInput) (*task.AccountUserListOutput, error) {
		tenantID := getTenantID(ctx)
		return service.SearchUsers(ctx, tenantID, input.Query)
	})
}
