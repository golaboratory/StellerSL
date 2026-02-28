package project

import "github.com/user/stellersl/backend/internal/api/task"

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

type ProjectIDInput struct {
	ID string `path:"id" example:"00000000-0000-0000-0000-000000000000"`
}

type ProjectUserAssignmentInput struct {
	ProjectIDInput
	Body struct {
		UserID string `json:"user_id" minLength:"1"`
	}
}

type ProjectUserUnassignmentInput struct {
	ProjectIDInput
	UserID string `path:"user_id" example:"00000000-0000-0000-0000-000000000000"`
}

type ProjectTaskListOutput struct {
	Body struct {
		Items []task.TaskItem `json:"items"`
	}
}
