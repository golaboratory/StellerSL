-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: SearchUsers :many
SELECT * FROM users 
WHERE (email ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateUser :one
UPDATE users
SET name = $2, avatar_url = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (tenant_id, email, password_hash, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDashboardStats :one
SELECT 
    COUNT(*) as total_tasks,
    COUNT(*) FILTER (WHERE status != 'done') as pending_tasks,
    COUNT(*) FILTER (WHERE status = 'done') as completed_tasks
FROM tasks 
WHERE deleted_at IS NULL AND (assigned_to = $1 OR assigned_to IS NULL);

-- name: GetDailyActivity :many
SELECT 
    logged_at::text as date,
    COUNT(*) FILTER (WHERE action = 'task_created') as created_count,
    COUNT(*) FILTER (WHERE action = 'task_completed') as completed_count
FROM activity_logs
WHERE user_id = $1 AND logged_at > CURRENT_DATE - INTERVAL '7 days'
GROUP BY logged_at
ORDER BY logged_at ASC;

-- name: GetRecentActivity :many
SELECT l.id, l.action, l.logged_at::text as date, COALESCE(t.title, 'Unknown Task')::text as task_title
FROM activity_logs l
LEFT JOIN tasks t ON l.task_id = t.id
WHERE l.user_id = $1
ORDER BY l.id DESC
LIMIT 10;

-- name: GetUserGrowth :one
SELECT * FROM user_growth WHERE user_id = $1;

-- name: CreateUserGrowth :exec
INSERT INTO user_growth (user_id) VALUES ($1) ON CONFLICT DO NOTHING;

-- name: AddExp :exec
UPDATE user_growth SET exp = exp + $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1;

-- name: UpdateLevel :exec
UPDATE user_growth SET level = (exp / 100) + 1 WHERE user_id = $1;

-- name: ListUserBadges :many
SELECT b.*
FROM badges b
JOIN user_badges ub ON b.id = ub.badge_id
WHERE ub.user_id = $1;

-- name: ListAllBadges :many
SELECT * FROM badges;

-- name: AwardBadge :exec
INSERT INTO user_badges (user_id, badge_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: GetTenantByDomain :one
SELECT * FROM tenants WHERE domain = $1;

-- name: ListProjects :many
SELECT * FROM projects 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateProject :one
INSERT INTO projects (tenant_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: UpdateProject :one
UPDATE projects SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: AssignProjectUser :exec
INSERT INTO project_users (project_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: UnassignProjectUser :exec
DELETE FROM project_users WHERE project_id = $1 AND user_id = $2;

-- name: ListProjectMembers :many
SELECT u.*
FROM users u
JOIN project_users pu ON u.id = pu.user_id
WHERE pu.project_id = $1;

-- name: ListTasks :many
SELECT * FROM tasks 
WHERE deleted_at IS NULL AND project_id = COALESCE($1, project_id) 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteTask :exec
UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (tenant_id, project_id, assigned_to, title, description, status, priority, due_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks 
SET title = $2, description = $3, status = $4, priority = $5, due_date = $6, updated_at = CURRENT_TIMESTAMP 
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateTaskStatus :one
UPDATE tasks SET status = $2, updated_at = CURRENT_TIMESTAMP 
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: BulkUpdateTasksStatus :exec
UPDATE tasks SET status = $2, updated_at = CURRENT_TIMESTAMP 
WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL;

-- name: BulkDeleteTasks :exec
UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP 
WHERE id = ANY($1::uuid[]);

-- name: CreateTeam :one
INSERT INTO teams (tenant_id, name) VALUES ($1, $2) RETURNING *;

-- name: ListTeams :many
SELECT * FROM teams 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: AddTeamMember :exec
INSERT INTO team_members (team_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (team_id, user_id) DO UPDATE SET role = EXCLUDED.role;

-- name: RemoveTeamMember :exec
DELETE FROM team_members WHERE team_id = $1 AND user_id = $2;

-- name: ListTeamMembers :many
SELECT u.*
FROM users u
JOIN team_members tm ON u.id = tm.user_id
WHERE tm.team_id = $1;
