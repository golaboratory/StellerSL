package dashboard

type DashboardOutput struct {
	Body struct {
		TotalTasks     int64 `json:"total_tasks"`
		PendingTasks   int64 `json:"pending_tasks"`
		CompletedTasks int64 `json:"completed_tasks"`
		DailyActivity  []struct {
			Date      string `json:"date"`
			Created   int64  `json:"created"`
			Completed int64  `json:"completed"`
		} `json:"daily_activity"`
	}
}
