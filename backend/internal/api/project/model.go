package project

type ProjectInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1"`
		Description string `json:"description"`
	}
}

type ProjectOutput struct {
	Body struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
}

type ProjectItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectListOutput struct {
	Body struct {
		Items []ProjectItem `json:"items"`
	}
}
