package auth

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
