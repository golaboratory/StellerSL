-- Add administrator flag to users (POST /auth/reset-password is admin-only)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE;

-- Add a full timestamp to activity_logs (logged_at is DATE only).
-- Needed for hour-window badges: lunch_hustler, golden_hour, hat_trick.
ALTER TABLE activity_logs ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
