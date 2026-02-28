-- name: GetUserByEmail :one
SELECT * FROM users WHERE tenant_id = $1 AND email = $2;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetDashboardStats :one
SELECT 
    COUNT(*) as total_tasks,
    COUNT(*) FILTER (WHERE status != 'done') as pending_tasks,
    COUNT(*) FILTER (WHERE status = 'done') as completed_tasks
FROM tasks 
WHERE tenant_id = $1 AND (assigned_to = $2 OR assigned_to IS NULL);

-- name: GetDailyActivity :many
SELECT 
    logged_at,
    COUNT(*) FILTER (WHERE action = 'task_created') as created_count,
    COUNT(*) FILTER (WHERE action = 'task_completed') as completed_count
FROM activity_logs
WHERE tenant_id = $1 AND user_id = $2 AND logged_at > CURRENT_DATE - INTERVAL '7 days'
GROUP BY logged_at
ORDER BY logged_at ASC;

-- name: GetTenantByDomain :one
SELECT * FROM tenants WHERE domain = $1;

-- name: ListProjects :many
SELECT * FROM projects WHERE tenant_id = $1 ORDER BY created_at DESC;

-- name: CreateProject :one
INSERT INTO projects (tenant_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTasks :many
SELECT * FROM tasks WHERE tenant_id = $1 AND project_id = COALESCE($2, project_id) ORDER BY created_at DESC;

-- name: CreateTask :one
INSERT INTO tasks (tenant_id, project_id, title, description, status, due_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateTaskStatus :one
UPDATE tasks SET status = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1 AND tenant_id = $2;
