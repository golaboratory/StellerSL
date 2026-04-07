-- Notifications table for in-app notification center
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL, -- 'badge_earned', 'level_up', 'team_invite', 'task_assigned'
    title TEXT NOT NULL,
    message TEXT,
    data JSONB, -- flexible payload (badge_id, task_id, etc.)
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id, read_at);
CREATE INDEX idx_notifications_tenant ON notifications(tenant_id);

-- RLS policy for notifications
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_notifications_policy ON notifications
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);
