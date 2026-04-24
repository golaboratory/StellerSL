# COMPLETE.md

実装完了済みタスクの記録。

---

## Phase 1: クリティカルな不具合修正（セキュリティ・データ整合性） — 全完了

### 1-1. JWT認証ミドルウェアの実装
- **対応内容**: Chi ミドルウェアとしてJWT検証を実装。`/auth/login`, `/auth/register`, `/openapi.json`, `/docs`, `/uploads/*` はスキップ。`api.go` のデフォルトユーザーIDフォールバックを401エラーに変更
- **対象ファイル**: `backend/cmd/main.go`, `backend/internal/api/api.go`

### 1-2. SQLインジェクション脆弱性の修正
- **対応内容**: `WithTenant` 内で `uuid.Parse(tenantID)` によるバリデーションを追加。不正なテナントIDはエラーで即座に拒否
- **対象ファイル**: `backend/internal/db/custom.go`

### 1-3. 平文パスワードフォールバックの削除
- **対応内容**: `auth/service.go` の `Login()` から `else if` 分岐（平文比較）を完全に削除。bcrypt検証のみに統一
- **対象ファイル**: `backend/internal/api/auth/service.go`

### 1-4. CORSミドルウェアの設定 (旧 Phase 8: 8-1)
- **対応内容**: `github.com/go-chi/cors` を追加。環境変数 `CORS_ORIGINS` で許可オリジンを制御（デフォルト: `http://localhost:5173`）
- **対象ファイル**: `backend/cmd/main.go`

---

## Phase 2: バッジシステムの完成 — 全完了（シークレットバッジを除く）

### 2-1. バッジシードデータの作成
- **対応内容**: `_initdb/07_badges.sql` を作成。`_memo/バッジの仕様.md` に定義された全19種のバッジをINSERT
- **対象バッジ**: first_task, three_strike, rookie_planner, daily_runner_bronze/silver/gold, comeback_hero, weekend_master, task_crusher_1〜5, early_bird, midnight_owl, quick_strike, all_nighter, organizer, multi_tasker

### 2-2. バッジエンジンの未実装ロジック追加
- **対応内容**: `engine.go` を全面的に拡張。全19種のバッジ判定ロジックを実装。`AwardBadges` シグネチャに `taskID ...uuid.UUID` オプション引数を追加（quick_strike判定用）
- **対象ファイル**: `backend/internal/api/gamification/engine.go`, `backend/sql/query.sql`

### 2-3. シークレットバッジの実装
- **決定事項**: Q1=D — Phase 2では対応しない。通常バッジの完成を優先し、後続フェーズに回す

### 2-4. ストリーク追跡のためのスキーマ拡張
- **決定事項**: Q2=C — 専用テーブル `user_streaks` を新設
- **対応内容**: `user_streaks` テーブル（user_id, streak_type, current_count, max_count, last_date）を作成。RLSポリシー追加。`StreakTracker` 構造体を実装し、タスク完了時に自動更新
- **対象ファイル**: `backend/sql/migrations/000004_user_streaks.up.sql`, `backend/sql/schema.sql`, `backend/internal/api/gamification/streak.go`

### 2-5. キャラクター自動進化 (旧 Phase 9: 9-3)
- **決定事項**: Q=A — レベル到達で自動進化
- **対応内容**: Lv.1-4→"egg", Lv.5-9→"chick", Lv.10-19→"chicken", Lv.20+→"phoenix"。タスク完了時に `UpdateLevel` 後に自動判定
- **対象ファイル**: `backend/internal/api/gamification/character.go`

---

## Phase 3: アクティビティログの実装 — 全完了

### 3-1. アクティビティログ書き込みの実装
- **対応内容**: タスクの Create/BulkCreate 時に `task_created`、Update/UpdateStatus/BulkUpdateStatus で status=="done" 時に `task_completed` を `activity_logs` にINSERT。`Create()` と `BulkCreate()` のシグネチャに `userID` パラメータを追加
- **対象ファイル**: `backend/internal/api/task/service.go`, `backend/internal/api/task/handler.go`, `backend/sql/query.sql`

### 3-2. ダッシュボードのアクティビティ表示修正
- **対応内容**: 3-1完了により `activity_logs` にデータが蓄積されるようになり、`GetRecentActivity` / `GetDailyActivity` クエリが正常に動作

---

## Phase 4: 不足しているAPIエンドポイント — 全完了

### 4-1. 単一タスク取得エンドポイント (GET /tasks/{id})
- **対応内容**: `task/handler.go` に `get-task` オペレーション追加、`task/service.go` に `Get()` メソッド追加
- **対象ファイル**: `backend/internal/api/task/handler.go`, `backend/internal/api/task/service.go`

### 4-2. チーム更新・削除エンドポイント
- **対応内容**: `PUT /teams/{id}` (update-team), `DELETE /teams/{id}` (delete-team) を追加。削除はownerのみ許可
- **対象ファイル**: `backend/internal/api/team/handler.go`, `backend/internal/api/team/service.go`, `backend/sql/query.sql`

### 4-3. 単一チーム取得エンドポイント (GET /teams/{id})
- **対応内容**: `GET /teams/{id}` (get-team) エンドポイント追加
- **対象ファイル**: `backend/internal/api/team/handler.go`, `backend/internal/api/team/service.go`, `backend/sql/query.sql`

### 4-4. 現在のユーザー情報取得エンドポイント (GET /auth/me)
- **対応内容**: `GET /auth/me` (get-me) エンドポイント追加。JWTからユーザーIDを取得し、ユーザー情報を返却
- **対象ファイル**: `backend/internal/api/auth/handler.go`, `backend/internal/api/auth/service.go`, `backend/internal/api/auth/model.go`

---

## Phase 5: タスク機能の拡充（部分完了）

### 5-3. タスクのステータス遷移ルールの実装
- **決定事項**: Q3=C — 制限付き自由遷移。done→他ステータスに戻す場合はEXP減算とバッジ再評価
- **対応内容**: `Update()` と `UpdateStatus()` で更新前のステータスを取得。done→非done 遷移時に `SubtractExp(10)` + `UpdateLevel` + `UpdateCharacterType` + `ReEvaluateBadges` + `InsertActivityLog(task_uncompleted)`。同一ステータスへの遷移では重複EXP付与を防止。時間系バッジ(early_bird等)は永続、回数系バッジは再評価で剥奪
- **対象ファイル**: `backend/internal/api/task/service.go`, `backend/internal/api/gamification/engine.go`

---

## Phase 6: フロントエンド品質改善（部分完了）

### 6-2. APIクライアントの統一
- **対応内容**: `Gantt.vue`, `CalendarView.vue`, `Signup.vue` を `Configuration` + `basePath` から `axiosInstance` パターンに統一
- **対象ファイル**: `frontend/src/views/Gantt.vue`, `frontend/src/views/CalendarView.vue`, `frontend/src/views/Signup.vue`

### 6-5. 不要コードの削除
- **対応内容**: 未使用の `HelloWorld.vue` を削除
- **対象ファイル**: `frontend/src/components/HelloWorld.vue`（削除）

---

## Phase 8: インフラ・運用改善（部分完了）

### 8-4. リフレッシュトークンの実装
- **決定事項**: Q5=D — 現状維持（24時間固定）。対応不要

---

## Phase 9: 追加機能（部分完了）

### 9-1. パスワードリセット機能
- **決定事項**: B+C — 管理者手動リセット + ログイン中パスワード変更
- **対応内容**: `PUT /auth/password` (change-password) と `POST /auth/reset-password` (reset-password) エンドポイントを追加。bcrypt検証→新ハッシュ生成→DB更新
- **対象ファイル**: `backend/internal/api/auth/handler.go`, `backend/internal/api/auth/service.go`, `backend/internal/api/auth/model.go`

### 9-2. チームの権限管理
- **決定事項**: A — 3階層ロールモデル (owner/admin/member)
- **対応内容**: `requireRole()` ヘルパーを実装。AddMember/RemoveMember にowner/admin権限チェック追加。DELETE /teams はowner限定。owner削除の防止ロジック追加
- **対象ファイル**: `backend/internal/api/team/service.go`, `backend/internal/api/team/handler.go`

### 9-3. キャラクタータイプの拡張
- → Phase 2: 2-5 として実装済み

### 9-4. 通知機能（バックエンドAPI実装済み）
- **決定事項**: B — アプリ内通知センター + ブラウザプッシュ通知
- **対応内容**: `notifications` テーブル新設。通知CRUDエンドポイント4本（GET /notifications, GET /notifications/unread/count, PATCH /notifications/{id}/read, PATCH /notifications/read-all）を実装。フロントエンドのUI実装は未完了
- **対象ファイル**: `backend/sql/migrations/000005_notifications.up.sql`, `backend/sql/schema.sql`, `backend/internal/api/notification/handler.go`, `backend/internal/api/notification/service.go`, `backend/internal/api/notification/model.go`, `backend/internal/api/api.go`

---

## 仕様確認 — 全回答済み

| # | 項目 | 回答 |
|---|---|---|
| Q1 | シークレットバッジ | D: 後続フェーズに回す |
| Q2 | ストリーク追跡方式 | C: 専用テーブル `user_streaks` 新設 |
| Q3 | ステータス遷移ルール | C: 制限付き自由遷移（done→他でEXP減算+バッジ再評価） |
| Q4 | 本番デプロイ方針 | B: VPS + Docker Compose + Let's Encrypt |
| Q5 | リフレッシュトークン | D: 現状維持（24時間固定） |
| 9-1 | パスワードリセット | B+C: 管理者手動リセット + ログイン中パスワード変更 |
| 9-2 | チーム権限管理 | A: 3階層ロールモデル |
| 9-3 | キャラクタータイプ | A: レベル到達で自動進化 |
| 9-4 | 通知機能 | B: アプリ内通知センター + ブラウザプッシュ通知 |
