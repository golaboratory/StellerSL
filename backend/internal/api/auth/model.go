package auth

import "github.com/danielgtaylor/huma/v2"

type RegisterInput struct {
	Body struct {
		TenantID     string `json:"tenant_id" required:"true"`
		Email        string `json:"email" format:"email" required:"true"`
		Password     string `json:"password" minLength:"8" required:"true"`
		Name         string `json:"name" required:"true"`
		InviteTeamID string `json:"invite_team_id,omitempty"`
	}
}

type LoginInput struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Password string `json:"password" minLength:"8"`
	}
}

type LoginOutput struct {
	Body struct {
		Token string `json:"token"`
		User  struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
	}
}

type UpdateProfileInput struct {
	Body struct {
		Name      string `json:"name" minLength:"1"`
		AvatarUrl string `json:"avatar_url"`
	}
}

type AvatarUploadInput struct {
	File huma.FormFile `multipart:"file"`
}

type AvatarUploadOutput struct {
	Body struct {
		Url string `json:"url"`
	}
}

type MeOutput struct {
	Body struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarUrl string `json:"avatar_url"`
		IsAdmin   bool   `json:"is_admin"`
	}
}

type ChangePasswordInput struct {
	Body struct {
		CurrentPassword string `json:"current_password" minLength:"8"`
		NewPassword     string `json:"new_password" minLength:"8"`
	}
}

type ResetPasswordInput struct {
	Body struct {
		UserID      string `json:"user_id"`
		NewPassword string `json:"new_password" minLength:"8"`
	}
}
