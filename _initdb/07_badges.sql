-- Badge seed data (19 badges from バッジの仕様.md)

-- 1. はじまりのバッジ（導入・チュートリアル）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000001', 'ファーストステップ', '初めてタスクを1つ完了する', 'pi-check', 'first_task'),
('10000000-0000-0000-0000-000000000002', 'ルーキー・プランナー', '1日のうちにタスクを5つ登録する', 'pi-pencil', 'rookie_planner'),
('10000000-0000-0000-0000-000000000003', 'スリーストライク', '累計でタスクを3つ完了する', 'pi-bolt', 'three_strike');

-- 2. コンボ・継続バッジ（習慣化）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000004', 'デイリーランナー（ブロンズ）', '3日連続でタスクを完了する', 'pi-calendar', 'daily_runner_bronze'),
('10000000-0000-0000-0000-000000000005', 'デイリーランナー（シルバー）', '7日連続でタスクを完了する', 'pi-calendar-plus', 'daily_runner_silver'),
('10000000-0000-0000-0000-000000000006', 'デイリーランナー（ゴールド）', '30日連続でタスクを完了する', 'pi-trophy', 'daily_runner_gold'),
('10000000-0000-0000-0000-000000000007', 'カムバックヒーロー', '7日以上の空白から復帰してタスクを完了する', 'pi-replay', 'comeback_hero'),
('10000000-0000-0000-0000-000000000008', 'ウィークエンドマスター', '土日にタスクを合計10個完了する', 'pi-sun', 'weekend_master');

-- 3. マイルストーンバッジ（累積・やり込み）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000009', 'タスククラッシャー Lv.1', 'タスクを累計10個完了する', 'pi-star', 'task_crusher_1'),
('10000000-0000-0000-0000-000000000010', 'タスククラッシャー Lv.2', 'タスクを累計50個完了する', 'pi-star-fill', 'task_crusher_2'),
('10000000-0000-0000-0000-000000000011', 'タスククラッシャー Lv.3', 'タスクを累計100個完了する', 'pi-verified', 'task_crusher_3'),
('10000000-0000-0000-0000-000000000012', 'タスククラッシャー Lv.4', 'タスクを累計500個完了する', 'pi-crown', 'task_crusher_4'),
('10000000-0000-0000-0000-000000000013', 'タスククラッシャー Lv.5', 'タスクを累計1000個完了する', 'pi-shield', 'task_crusher_5');

-- 4. スピード＆タイムバッジ（プレイスタイル・遊び心）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000014', 'アーリーバード', '午前4:00〜8:59の間にタスクを完了する', 'pi-sun', 'early_bird'),
('10000000-0000-0000-0000-000000000015', 'ミッドナイトオウル', '午前0:00〜3:59の間にタスクを完了する', 'pi-moon', 'midnight_owl'),
('10000000-0000-0000-0000-000000000016', 'クイックストライク', 'タスクを登録してから1時間以内に完了する', 'pi-stopwatch', 'quick_strike'),
('10000000-0000-0000-0000-000000000017', '一網打尽', '1日のうちにタスクを10個完了する', 'pi-bolt', 'all_nighter');

-- 5. プロジェクト・整理バッジ（機能活用）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000018', 'オーガナイザー', 'プロジェクトを3つ作成する', 'pi-folder', 'organizer'),
('10000000-0000-0000-0000-000000000019', 'マルチタスカー', '3つの異なるプロジェクトのタスクを同日に完了する', 'pi-th-large', 'multi_tasker');

-- === Phase B（次期実装推奨 — 10種） ===

-- 1. はじまりのバッジ（補完）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000020', 'はじめの一歩', 'プロジェクトを初めて1つ作成する', 'pi-folder-open', 'first_project'),
('10000000-0000-0000-0000-000000000021', 'チームプレイヤー', 'チームに初めて参加する', 'pi-users', 'first_team_join');

-- 2. コンボ・継続バッジ（上位・曜日）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000022', 'デイリーランナー（プラチナ）', '累計最長ストリーク100日以上を達成する', 'pi-star', 'daily_runner_platinum'),
('10000000-0000-0000-0000-000000000023', 'マンデーモチベーター', '月曜日にタスクを5つ完了する', 'pi-flag', 'monday_motivator'),
('10000000-0000-0000-0000-000000000024', 'フライデーフィニッシャー', '金曜日にタスクを5つ完了して週を締める', 'pi-thumbs-up', 'friday_finisher');

-- 4. スピード＆タイムバッジ（拡充）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000025', '電光石火', 'タスクを登録してから10分以内に完了する', 'pi-forward', 'lightning_fast'),
('10000000-0000-0000-0000-000000000026', 'ランチタイムハッスル', '12:00〜12:59の間にタスクを3つ完了する', 'pi-clock', 'lunch_hustler'),
('10000000-0000-0000-0000-000000000027', 'ゴールデンアワー', '17:00〜17:59の間にタスクを3つ完了する', 'pi-palette', 'golden_hour'),
('10000000-0000-0000-0000-000000000028', 'ハットトリック', '1時間以内にタスクを3つ完了する', 'pi-chart-line', 'hat_trick');

-- 5. プロジェクト・整理バッジ（実用）
INSERT INTO badges (id, name, description, icon_slug, requirement_type) VALUES
('10000000-0000-0000-0000-000000000029', 'ゼロインボックス', 'プロジェクト内の全タスクを完了状態にする', 'pi-inbox', 'zero_inbox');
