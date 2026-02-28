package gamification

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	UserID string
}

func RegisterHandlers(api huma.API, service *Service, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user-growth",
		Method:      http.MethodGet,
		Path:        "/growth",
		Summary:     "User Growth",
	}, func(ctx context.Context, input *struct{}) (*GrowthOutput, error) {
		auth, _ := getAuth(ctx)
		return service.GetGrowth(ctx, auth.UserID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-badges",
		Method:      http.MethodGet,
		Path:        "/badges",
		Summary:     "User Badges",
	}, func(ctx context.Context, input *struct{}) (*BadgeListOutput, error) {
		auth, _ := getAuth(ctx)
		return service.ListBadges(ctx, auth.UserID)
	})
}
