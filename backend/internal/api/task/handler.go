package task

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	TenantID string
	UserID   string
}

type ListTasksInput struct {
	Limit  int32 `query:"limit" default:"50" maximum:"200"`
	Offset int32 `query:"offset" default:"0"`
}

func RegisterHandlers(api huma.API, service *Service, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "list-tasks",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "List Tasks",
	}, func(ctx context.Context, input *ListTasksInput) (*TaskListOutput, error) {
		auth, _ := getAuth(ctx)
		return service.List(ctx, auth.TenantID, input.Limit, input.Offset)
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-task",
		Method:      http.MethodPost,
		Path:        "/tasks",
		Summary:     "Create Task",
	}, func(ctx context.Context, input *TaskInput) (*TaskOutput, error) {
		auth, _ := getAuth(ctx)
		return service.Create(ctx, auth.TenantID, *input)
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      http.MethodPut,
		Path:        "/tasks/{id}",
		Summary:     "Update Task",
	}, func(ctx context.Context, input *struct {
		ID string `path:"id"`
		TaskInput
	}) (*TaskOutput, error) {
		auth, _ := getAuth(ctx)
		return service.Update(ctx, auth.TenantID, auth.UserID, input.ID, input.TaskInput)
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-task-status",
		Method:      http.MethodPatch,
		Path:        "/tasks/{id}/status",
		Summary:     "Update Task Status",
	}, func(ctx context.Context, input *TaskStatusUpdateInput) (*TaskOutput, error) {
		auth, _ := getAuth(ctx)
		return service.UpdateStatus(ctx, auth.TenantID, auth.UserID, input.ID, input.Body.Status)
	})

	huma.Register(api, huma.Operation{
		OperationID: "bulk-create-tasks",
		Method:      http.MethodPost,
		Path:        "/tasks/bulk",
		Summary:     "Bulk Create Tasks",
	}, func(ctx context.Context, input *BulkTaskCreateInput) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.BulkCreate(ctx, auth.TenantID, *input); err != nil {
			return nil, huma.Error500InternalServerError("Bulk create failed")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "bulk-update-tasks-status",
		Method:      http.MethodPatch,
		Path:        "/tasks/bulk/status",
		Summary:     "Bulk Update Tasks Status",
	}, func(ctx context.Context, input *BulkTaskUpdateInput) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.BulkUpdateStatus(ctx, auth.TenantID, auth.UserID, input.Body.Status, input.Body.IDs); err != nil {
			return nil, huma.Error500InternalServerError("Bulk update failed")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "bulk-delete-tasks",
		Method:      http.MethodDelete,
		Path:        "/tasks/bulk",
		Summary:     "Bulk Delete Tasks",
	}, func(ctx context.Context, input *BulkTaskDeleteInput) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.BulkDelete(ctx, auth.TenantID, input.Body.IDs); err != nil {
			return nil, huma.Error500InternalServerError("Bulk delete failed")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-task",
		Method:      http.MethodDelete,
		Path:        "/tasks/{id}",
		Summary:     "Delete Task",
	}, func(ctx context.Context, input *struct{ ID string `path:"id"` }) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.Delete(ctx, auth.TenantID, input.ID); err != nil {
			return nil, huma.Error500InternalServerError("Delete failed")
		}
		return nil, nil
	})
}
