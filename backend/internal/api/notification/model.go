package notification

type NotificationItem struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Message   *string `json:"message,omitempty"`
	ReadAt    *string `json:"read_at,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type NotificationListOutput struct {
	Body struct {
		Items []NotificationItem `json:"items"`
	}
}

type UnreadCountOutput struct {
	Body struct {
		Count int32 `json:"count"`
	}
}

type AuthInfo struct {
	TenantID string
	UserID   string
}
