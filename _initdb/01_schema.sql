-- 1. Tenants (Root for logical separation)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    domain TEXT UNIQUE NOT NULL, -- Used for domain-based identification
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Users (Scoped by Tenant)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    avatar_url TEXT,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, email)
);

-- 3. Teams (For collaborative management)
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE team_members (
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member', -- owner, admin, member
    PRIMARY KEY (team_id, user_id)
);

-- 4. Projects (Task groups)
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Project Assignments (User access to projects)
CREATE TABLE project_users (
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, user_id)
);

-- 6. Tasks (The core unit)
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'todo', -- todo, doing, done
    priority INTEGER DEFAULT 0,
    due_date TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. Gamification: Character Growth
CREATE TABLE user_growth (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    level INTEGER NOT NULL DEFAULT 1,
    exp INTEGER NOT NULL DEFAULT 0,
    character_type TEXT NOT NULL DEFAULT 'default',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 8. Gamification: Badges
CREATE TABLE badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    icon_slug TEXT NOT NULL, -- pi-star, pi-check, etc (PrimeIcons)
    requirement_type TEXT NOT NULL -- task_count, login_streak, etc
);

CREATE TABLE user_badges (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    badge_id UUID REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, badge_id)
);

-- 9. Activity Logs (For Graphing daily counts)
CREATE TABLE activity_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    action TEXT NOT NULL, -- task_created, task_completed
    logged_at DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- time-of-day for hour-window badges
);

-- 10. User Streaks (For tracking consecutive activity)
CREATE TABLE user_streaks (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    streak_type TEXT NOT NULL, -- 'daily_task_completion'
    current_count INTEGER NOT NULL DEFAULT 0,
    max_count INTEGER NOT NULL DEFAULT 0,
    last_date DATE,
    PRIMARY KEY (user_id, streak_type)
);

-- 11. Notifications (In-app notification center)
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL, -- 'badge_earned', 'level_up', 'team_invite', 'task_assigned'
    title TEXT NOT NULL,
    message TEXT,
    data JSONB,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Multi-tenancy & Performance Indexes
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_teams_tenant_id ON teams(tenant_id);
CREATE INDEX idx_projects_tenant_id ON projects(tenant_id);
CREATE INDEX idx_tasks_tenant_id ON tasks(tenant_id);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_activity_logs_lookup ON activity_logs(tenant_id, user_id, logged_at);
CREATE INDEX idx_notifications_user ON notifications(user_id, read_at);
CREATE INDEX idx_notifications_tenant ON notifications(tenant_id);

-- Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE team_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE project_users ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_growth ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_badges ENABLE ROW LEVEL SECURITY;
ALTER TABLE activity_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_streaks ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_users_policy ON users
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_teams_policy ON teams
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_team_members_policy ON team_members
    USING (EXISTS (
        SELECT 1 FROM teams
        WHERE teams.id = team_members.team_id
        AND teams.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

CREATE POLICY tenant_projects_policy ON projects
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_project_users_policy ON project_users
    USING (EXISTS (
        SELECT 1 FROM projects
        WHERE projects.id = project_users.project_id
        AND projects.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

CREATE POLICY tenant_tasks_policy ON tasks
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_user_growth_policy ON user_growth
    USING (EXISTS (
        SELECT 1 FROM users
        WHERE users.id = user_growth.user_id
        AND users.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

CREATE POLICY tenant_user_badges_policy ON user_badges
    USING (EXISTS (
        SELECT 1 FROM users
        WHERE users.id = user_badges.user_id
        AND users.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

CREATE POLICY tenant_activity_logs_policy ON activity_logs
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_user_streaks_policy ON user_streaks
    USING (EXISTS (
        SELECT 1 FROM users
        WHERE users.id = user_streaks.user_id
        AND users.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

CREATE POLICY tenant_notifications_policy ON notifications
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);
