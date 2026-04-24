-- User Streaks table for tracking consecutive activity
CREATE TABLE IF NOT EXISTS user_streaks (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    streak_type TEXT NOT NULL, -- 'daily_task_completion'
    current_count INTEGER NOT NULL DEFAULT 0,
    max_count INTEGER NOT NULL DEFAULT 0,
    last_date DATE,
    PRIMARY KEY (user_id, streak_type)
);

-- RLS policy for user_streaks (via users table tenant check)
ALTER TABLE user_streaks ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_user_streaks_policy ON user_streaks
    USING (EXISTS (
        SELECT 1 FROM users
        WHERE users.id = user_streaks.user_id
        AND users.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));
