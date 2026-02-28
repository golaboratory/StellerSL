package team

type TeamInput struct {
	Body struct {
		Name string `json:"name" minLength:"1"`
	}
}

type TeamItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeamListOutput struct {
	Body struct {
		Items []TeamItem `json:"items"`
	}
}

type TeamMemberInput struct {
	ID   string `path:"id"`
	Body struct {
		UserID string `json:"user_id"`
		Role   string `json:"role" enum:"owner,admin,member" default:"member"`
	}
}

type TeamMemberUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type TeamMemberListOutput struct {
	Body struct {
		Items []TeamMemberUser `json:"items"`
	}
}
