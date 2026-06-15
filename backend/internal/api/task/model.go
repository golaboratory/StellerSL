package task

type TaskItem struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	AssignedTo  string `json:"assigned_to,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	DueDate     string `json:"due_date,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type TaskListOutput struct {
	Body struct {
		Items []TaskItem `json:"items"`
	}
}

type TaskInput struct {
	Body struct {
		ProjectID   string `json:"project_id,omitempty"`
		AssignedTo  string `json:"assigned_to,omitempty"`
		Title       string `json:"title" minLength:"1"`
		Description string `json:"description"`
		Status      string `json:"status" enum:"todo,doing,done" default:"todo"`
		Priority    int    `json:"priority" default:"0"`
		DueDate     string `json:"due_date,omitempty" format:"date-time"`
	}
}

type TaskOutput struct {
	Body struct {
		ID          string `json:"id"`
		ProjectID   string `json:"project_id"`
		AssignedTo  string `json:"assigned_to,omitempty"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    int    `json:"priority"`
		DueDate     string `json:"due_date"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
}

type TaskStatusUpdateInput struct {
	ID   string `path:"id"`
	Body struct {
		Status string `json:"status" enum:"todo,doing,done"`
	}
}

type BulkTaskCreateItem struct {
	ProjectID   string `json:"project_id"`
	Title       string `json:"title" minLength:"1"`
	Description string `json:"description"`
}

type BulkTaskCreateInput struct {
	Body struct {
		Tasks []BulkTaskCreateItem `json:"tasks"`
	}
}

type BulkTaskUpdateInput struct {
	Body struct {
		IDs    []string `json:"ids"`
		Status string   `json:"status" enum:"todo,doing,done"`
	}
}

type BulkTaskDeleteInput struct {
	Body struct {
		IDs []string `json:"ids"`
	}
}

type AccountUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AccountUserListOutput struct {
	Body struct {
		Items []AccountUser `json:"items"`
	}
}
