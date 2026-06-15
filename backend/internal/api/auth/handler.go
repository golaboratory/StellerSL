package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	TenantID string
	UserID   string
}

func RegisterHandlers(api huma.API, service *Service, getTenantID func(context.Context) string, getAuth func(context.Context) (AuthInfo, error)) {
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/auth/register",
		Summary:     "User Registration",
	}, func(ctx context.Context, input *RegisterInput) (*struct{}, error) {
		if err := service.Register(ctx, *input); err != nil {
			return nil, huma.Error500InternalServerError("Failed to create user: " + err.Error())
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
		fmt.Println("tenantID", tenantID)
		resp, err := service.Login(ctx, *input, tenantID)
		fmt.Println("resp", resp)
		fmt.Println("err", err)
		if err != nil {
			return nil, huma.Error401Unauthorized("Invalid credentials")
		}
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-me",
		Method:      http.MethodGet,
		Path:        "/auth/me",
		Summary:     "Get Current User",
	}, func(ctx context.Context, input *struct{}) (*MeOutput, error) {
		auth, err := getAuth(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}
		return service.GetMe(ctx, auth.TenantID, auth.UserID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "change-password",
		Method:      http.MethodPut,
		Path:        "/auth/password",
		Summary:     "Change Password",
	}, func(ctx context.Context, input *ChangePasswordInput) (*struct{}, error) {
		auth, err := getAuth(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}
		if err := service.ChangePassword(ctx, auth.TenantID, auth.UserID, *input); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "reset-password",
		Method:      http.MethodPost,
		Path:        "/auth/reset-password",
		Summary:     "Reset Password (Admin)",
	}, func(ctx context.Context, input *ResetPasswordInput) (*struct{}, error) {
		auth, err := getAuth(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}
		isAdmin, err := service.IsAdmin(ctx, auth.TenantID, auth.UserID)
		if err != nil || !isAdmin {
			return nil, huma.Error403Forbidden("Administrator privileges required")
		}
		if err := service.ResetPassword(ctx, auth.TenantID, *input); err != nil {
			return nil, huma.Error500InternalServerError("Failed to reset password")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-profile",
		Method:      http.MethodPut,
		Path:        "/auth/profile",
		Summary:     "Update User Profile",
	}, func(ctx context.Context, input *UpdateProfileInput) (*struct{}, error) {
		auth, err := getAuth(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}
		if err := service.UpdateProfile(ctx, auth.TenantID, auth.UserID, *input); err != nil {
			return nil, huma.Error500InternalServerError("Failed to update profile")
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "upload-avatar",
		Method:      http.MethodPost,
		Path:        "/auth/avatar",
		Summary:     "Upload Avatar Image",
	}, func(ctx context.Context, input *AvatarUploadInput) (*AvatarUploadOutput, error) {
		auth, err := getAuth(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}

		data, err := io.ReadAll(input.File.File)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to read uploaded file")
		}

		url, err := service.UpdateAvatar(ctx, auth.TenantID, auth.UserID, input.File.Filename, data)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to upload avatar")
		}

		resp := &AvatarUploadOutput{}
		resp.Body.Url = url
		return resp, nil
	})
}
