-- 1. Add deleted_at to tasks for soft delete
ALTER TABLE tasks ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

-- 2. Seed initial badges based on the specification
INSERT INTO badges (name, description, icon_slug, requirement_type) VALUES
('ファーストステップ', '初めてタスクを1つ完了する', 'pi-check-circle', 'first_task'),
('ルーキー・プランナー', '一度に（または1日のうちに）タスクを5つ登録する', 'pi-plus-circle', 'rookie_planner'),
('スリーストライク', '累計でタスクを3つ完了する', 'pi-star', 'three_strike'),
('デイリーランナー (ブロンズ)', '3日連続でタスクを完了する', 'pi-calendar', 'daily_runner_bronze'),
('タスククラッシャー Lv.1', 'タスク累計 10 個完了', 'pi-bolt', 'task_crusher_1'),
('アーリーバード', '午前4:00〜午前8:59 の間にタスクを完了する', 'pi-sun', 'early_bird'),
('ミッドナイトオウル', '午前0:00〜午前3:59 の間にタスクを完了する', 'pi-moon', 'midnight_owl');
