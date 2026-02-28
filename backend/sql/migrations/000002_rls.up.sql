-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE team_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE project_users ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_growth ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_badges ENABLE ROW LEVEL SECURITY;
ALTER TABLE activity_logs ENABLE ROW LEVEL SECURITY;

-- We don't enable RLS on 'tenants' and 'badges' tables as they are global/lookup tables
-- OR we can enable them with public read access.

-- Define Policies based on 'app.current_tenant_id' session variable

-- 1. Users
CREATE POLICY tenant_users_policy ON users
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- 2. Teams
CREATE POLICY tenant_teams_policy ON teams
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- 3. Team Members (Join with teams to verify tenant)
CREATE POLICY tenant_team_members_policy ON team_members
    USING (EXISTS (
        SELECT 1 FROM teams WHERE teams.id = team_members.team_id 
        AND teams.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

-- 4. Projects
CREATE POLICY tenant_projects_policy ON projects
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- 5. Project Users (Join with projects)
CREATE POLICY tenant_project_users_policy ON project_users
    USING (EXISTS (
        SELECT 1 FROM projects WHERE projects.id = project_users.project_id 
        AND projects.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

-- 6. Tasks
CREATE POLICY tenant_tasks_policy ON tasks
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- 7. User Growth (Join with users)
CREATE POLICY tenant_user_growth_policy ON user_growth
    USING (EXISTS (
        SELECT 1 FROM users WHERE users.id = user_growth.user_id 
        AND users.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

-- 8. User Badges (Join with users)
CREATE POLICY tenant_user_badges_policy ON user_badges
    USING (EXISTS (
        SELECT 1 FROM users WHERE users.id = user_badges.user_id 
        AND users.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));

-- 9. Activity Logs
CREATE POLICY tenant_activity_logs_policy ON activity_logs
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);
