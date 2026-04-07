package notification

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ListInput struct {
	Limit  int32 `query:"limit" default:"50" maximum:"200"`
	Offset int32 `query:"offset" default:"0"`
}

func RegisterHandlers(api huma.API, service *Service, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "list-notifications",
		Method:      http.MethodGet,
		Path:        "/notifications",
		Summary:     "List Notifications",
	}, func(ctx context.Context, input *ListInput) (*NotificationListOutput, error) {
		auth, _ := getAuth(ctx)
		return service.List(ctx, auth.TenantID, auth.UserID, input.Limit, input.Offset)
	})

	huma.Register(api, huma.Operation{
		OperationID: "count-unread-notifications",
		Method:      http.MethodGet,
		Path:        "/notifications/unread/count",
		Summary:     "Count Unread Notifications",
	}, func(ctx context.Context, input *struct{}) (*UnreadCountOutput, error) {
		auth, _ := getAuth(ctx)
		return service.UnreadCount(ctx, auth.TenantID, auth.UserID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "mark-notification-read",
		Method:      http.MethodPatch,
		Path:        "/notifications/{id}/read",
		Summary:     "Mark Notification Read",
	}, func(ctx context.Context, input *struct{ ID string `path:"id"` }) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.MarkRead(ctx, auth.TenantID, auth.UserID, input.ID); err != nil {
			return nil, huma.Error500InternalServerError("Failed to mark notification")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "mark-all-notifications-read",
		Method:      http.MethodPatch,
		Path:        "/notifications/read-all",
		Summary:     "Mark All Notifications Read",
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
		auth, _ := getAuth(ctx)
		if err := service.MarkAllRead(ctx, auth.TenantID, auth.UserID); err != nil {
			return nil, huma.Error500InternalServerError("Failed to mark all read")
		}
		return nil, nil
	})
}
