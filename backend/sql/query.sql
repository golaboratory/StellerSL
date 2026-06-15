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
WHERE deleted_at IS NULL
  AND (sqlc.narg('project_id')::uuid IS NULL OR project_id = sqlc.narg('project_id'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('priority')::int IS NULL OR priority = sqlc.narg('priority'))
  AND (sqlc.narg('due_date_from')::timestamptz IS NULL OR due_date >= sqlc.narg('due_date_from'))
  AND (sqlc.narg('due_date_to')::timestamptz IS NULL OR due_date <= sqlc.narg('due_date_to'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

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
SET title = $2, description = $3, status = $4, priority = $5, due_date = $6, assigned_to = $7, updated_at = CURRENT_TIMESTAMP
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
SELECT u.id, u.email, u.name, u.avatar_url, tm.role
FROM users u
JOIN team_members tm ON u.id = tm.user_id
WHERE tm.team_id = $1
ORDER BY CASE tm.role WHEN 'owner' THEN 0 WHEN 'admin' THEN 1 ELSE 2 END, u.name ASC;

-- === Activity Log Queries ===

-- name: InsertActivityLog :exec
INSERT INTO activity_logs (tenant_id, user_id, task_id, action)
VALUES ($1, $2, $3, $4);

-- name: CountTasksCreatedToday :one
SELECT COUNT(*)::int as count FROM activity_logs
WHERE user_id = $1 AND action = 'task_created' AND logged_at = CURRENT_DATE;

-- name: CountTasksCompletedToday :one
SELECT COUNT(*)::int as count FROM activity_logs
WHERE user_id = $1 AND action = 'task_completed' AND logged_at = CURRENT_DATE;

-- name: CountTasksCompletedOnWeekends :one
SELECT COUNT(*)::int as count FROM activity_logs
WHERE user_id = $1 AND action = 'task_completed' AND EXTRACT(DOW FROM logged_at) IN (0, 6);

-- name: CountDistinctProjectsCompletedToday :one
SELECT COUNT(DISTINCT t.project_id)::int as count
FROM activity_logs al
JOIN tasks t ON al.task_id = t.id
WHERE al.user_id = $1 AND al.action = 'task_completed' AND al.logged_at = CURRENT_DATE;

-- name: GetLastActivityDate :one
-- Last completion date BEFORE today (activity logs are written before badge
-- evaluation, so today's completion must be excluded for comeback detection).
SELECT COALESCE(MAX(logged_at), '1970-01-01'::date)::date as last_date FROM activity_logs
WHERE user_id = $1 AND action = 'task_completed' AND logged_at < CURRENT_DATE;

-- === Streak Queries ===

-- name: GetStreak :one
SELECT * FROM user_streaks WHERE user_id = $1 AND streak_type = $2;

-- name: UpsertStreak :exec
INSERT INTO user_streaks (user_id, streak_type, current_count, max_count, last_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, streak_type) DO UPDATE SET
    current_count = EXCLUDED.current_count,
    max_count = GREATEST(user_streaks.max_count, EXCLUDED.max_count),
    last_date = EXCLUDED.last_date;

-- === Gamification Extra Queries ===

-- name: SubtractExp :exec
UPDATE user_growth SET exp = GREATEST(exp - $2, 0), updated_at = CURRENT_TIMESTAMP WHERE user_id = $1;

-- name: RemoveBadge :exec
DELETE FROM user_badges WHERE user_id = $1 AND badge_id = $2;

-- name: UpdateCharacterType :exec
UPDATE user_growth SET character_type = $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1;

-- name: CountProjectsByUser :one
SELECT COUNT(*)::int as count FROM project_users WHERE user_id = $1;

-- name: CountTeamMembershipsByUser :one
SELECT COUNT(*)::int as count FROM team_members WHERE user_id = $1;

-- name: CountTasksCompletedTodayInHour :one
-- Completions today within the given hour (server-local), e.g. 12 = 12:00-12:59.
SELECT COUNT(*)::int as count FROM activity_logs
WHERE user_id = $1 AND action = 'task_completed' AND logged_at = CURRENT_DATE
  AND EXTRACT(HOUR FROM created_at)::int = $2::int;

-- name: CountTasksCompletedLastHour :one
SELECT COUNT(*)::int as count FROM activity_logs
WHERE user_id = $1 AND action = 'task_completed' AND created_at >= NOW() - INTERVAL '1 hour';

-- name: GetProjectTaskCounts :one
SELECT COUNT(*)::int AS total, (COUNT(*) FILTER (WHERE status != 'done'))::int AS remaining
FROM tasks WHERE project_id = $1 AND deleted_at IS NULL;

-- === Team Extra Queries ===

-- name: GetTeam :one
SELECT * FROM teams WHERE id = $1;

-- name: UpdateTeam :one
UPDATE teams SET name = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;

-- name: GetTeamMemberRole :one
SELECT role FROM team_members WHERE team_id = $1 AND user_id = $2;

-- === Auth Extra Queries ===

-- name: UpdatePassword :exec
UPDATE users SET password_hash = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- === Notification Queries ===

-- name: CreateNotification :exec
INSERT INTO notifications (tenant_id, user_id, type, title, message, data)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListNotifications :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUnreadNotifications :many
SELECT * FROM notifications
WHERE user_id = $1 AND read_at IS NULL
ORDER BY created_at DESC
LIMIT $2;

-- name: CountUnreadNotifications :one
SELECT COUNT(*)::int as count FROM notifications
WHERE user_id = $1 AND read_at IS NULL;

-- name: MarkNotificationRead :exec
UPDATE notifications SET read_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications SET read_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND read_at IS NULL;
