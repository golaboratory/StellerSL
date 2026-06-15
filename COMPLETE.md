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

---

## 2026-06-11 実装スプリント（手順1〜5 一括実行） — 全完了

### 修正: reset-password の管理者権限チェック（旧TODOコメント解消）
- **対応内容**: `users.is_admin` カラムを追加（schema.sql / _initdb / migrations/000006）。`POST /auth/reset-password` は管理者のみ403で制限。`GET /auth/me` に `is_admin` を追加。シードユーザー (test@example.com) を管理者に設定
- **対象ファイル**: `backend/internal/api/auth/handler.go`, `service.go`, `model.go`

### 5-1. タスクのフィルタリング機能
- **対応内容**: `GET /tasks` に `status` / `priority` / `project_id` / `due_date_from` / `due_date_to` クエリパラメータを追加（sqlc.narg によるNULL許容フィルタ）。既存バグ修正: project_id が NULL のタスクが一覧から除外されていた COALESCE 比較を撤廃
- **対象ファイル**: `backend/sql/query.sql`, `backend/internal/api/task/handler.go`, `service.go`

### 5-2. タスク担当者の変更機能
- **対応内容**: `UpdateTask` クエリと `Update()` に `assigned_to` を追加。TaskItem / TaskOutput に `assigned_to` / `description` / `created_at` を追加し、`due_date` を RFC3339（NULL時は省略）に統一
- **対象ファイル**: `backend/sql/query.sql`, `backend/internal/api/task/*.go`

### 6-1. フォームバリデーションの強化
- **対応内容**: `src/lib/validation.ts`（メール形式・パスワード強度）を新設。Login / Signup / Tasks編集 / CalendarView / Projects / TeamDetail にフィールド単位のエラー表示と `:invalid` を追加
- **対象ファイル**: `frontend/src/lib/validation.ts`, 各ビュー

### 6-3 / 7-FE-6. ガントチャートの編集機能
- **対応内容**: バーのドラッグで `PUT /tasks/{id}` を呼び due_date を保存（全フィールド再送で上書き消失を防止）。モック日付を廃止し実データ（created_at〜due_date）で描画
- **対象ファイル**: `frontend/src/views/Gantt.vue`

### 6-4 / 7-FE-7. カレンダービューからのタスク作成・編集
- **対応内容**: 日付選択 +「Add Task」で選択日を期日にしたタスク作成ダイアログ、タスククリックで編集ダイアログを追加。プロジェクト選択 (Select) 対応
- **対象ファイル**: `frontend/src/views/CalendarView.vue`

### 6-6 / 7-FE-1 / 7-FE-2. チームUX改善
- **対応内容**: UUID手入力を廃止しユーザー検索ダイアログ（searchUsers）に変更。ロール表示を実データ化し owner/admin はロール変更ドロップダウン・メンバー削除・チーム名変更・チーム削除（ownerのみ）を操作可能に。バックエンド: チーム作成者を owner として自動登録（従来は作成者がメンバー外で管理操作不能だった）、ListTeamMembers がロールを返却、owner のロール降格をupsert経由でも防止
- **対象ファイル**: `frontend/src/views/TeamDetail.vue`, `backend/internal/api/team/*.go`, `backend/sql/query.sql`

### 7-2. バッジエンジンのユニットテスト
- **対応内容**: `engine_test.go` を新規作成。`db.Querier` フェイクで全29種の授与判定・剥奪ルール・ストリーク更新・キャラ進化閾値をテーブルドリブンで検証。エンジン/ストリークはテナントTX内の `db.Querier` を引数で受ける構造にリファクタ（RLS適用下でも正しく動作）し、時刻は注入可能に
- **対象ファイル**: `backend/internal/api/gamification/engine{,_test}.go`, `streak.go`, `character.go`

### 7-FE-3. 通知ベル + ドロワー
- **対応内容**: `stores/notification.ts`（未読数30秒ポーリング）、`NotificationBell.vue`（OverlayBadge付き）、`NotificationDrawer.vue`（一覧・既読化・全既読）を新設し Dashboard ツールバーに設置
- **対象ファイル**: `frontend/src/components/Notification*.vue`, `frontend/src/stores/notification.ts`

### 7-FE-4. プロフィールにパスワード変更UI
- **対応内容**: 現在/新規/確認の3フィールド + 強度チェック + `PUT /auth/password` 呼び出しを Profile に追加
- **対象ファイル**: `frontend/src/views/Profile.vue`

### 7-FE-5. キャラクター進化表示
- **対応内容**: `character_type`（egg/chick/chicken/phoenix）ベースの表示に修正（従来はレベル閾値のみでフェニックス段階が欠落）。和名ラベル表示
- **対象ファイル**: `frontend/src/views/Dashboard.vue`

### バッジ Phase B（10種）の実装
- **対応内容**: `first_project` / `first_team_join` / `daily_runner_platinum` / `monday_motivator` / `friday_finisher` / `lightning_fast` / `lunch_hustler` / `golden_hour` / `hat_trick` / `zero_inbox` をエンジンに実装。シード追加（_initdb/07 + migrations/000007）。時間帯判定用に `activity_logs.created_at` (timestamptz) を追加（migrations/000006）。実装済みバッジは計29種
- **既存バグ修正**: `rookie_planner` が一切発火しなかった（task_created トリガーでの評価が未配線）/ `all_nighter` が仕様に反して剥奪されていた / `comeback_hero` が初回完了で誤授与 / ストリーク更新がバッジ評価の後に実行され当日の streak バッジが翌日までずれていた / `Truncate(24h)` によるタイムゾーン依存の日界ずれ
- **対象ファイル**: `backend/internal/api/gamification/engine.go`, `backend/internal/api/task/service.go`, `backend/sql/*`

### バッジSVGアセットの配線
- **対応内容**: `list-badges` に `requirement_type` を追加し、`src/assets/badges/*.svg`（47個）を requirement_type で解決して表示（無い場合は PrimeIcons にフォールバック）
- **対象ファイル**: `backend/internal/api/gamification/model.go`, `frontend/src/lib/badgeIcons.ts`, `frontend/src/views/Dashboard.vue`

### E2E環境の復旧と全テスト緑化
- **対応内容**: playwright.config.ts の webServer を有効化（vite自動起動）。E2E全滅の原因を解消: ①Loginボタンのラベル不一致（Sign In→Login）②サインアップ/タスク作成のPOST完了を待たずに遷移するレース ③テナント共有データとのタスク名衝突（タイムスタンプで一意化）④`<Toast />` の二重マウント（App.vue と Dashboard.vue）⑤一括削除前の再選択漏れ。`@playwright/test` を 1.60 に更新（Ubuntu 26.04 対応）。vitest が E2E スペックを誤って拾う問題を exclude で解消。常時失敗していた `register_debug_test.go`（要DB・固定メール）を削除
- **検証結果**: バックエンド `go test ./...` 全パス / フロント vitest 6/6 / `pnpm build` (vue-tsc) 成功 / Playwright E2E **5/5 パス**

### ビルド健全化（既存の型エラー一掃）
- **対応内容**: `pnpm build` が生成クライアント由来の型エラー等で失敗していた問題を解消（updateTask/updateTeam のパラメータ名不一致、type-only import、未使用変数、`auth.tenantID` 不存在参照、v-calendar/frappe-gantt の型、tsconfig の erasableSyntaxOnly と生成コードの衝突）

---

## 2026-06-12 レスポンシブ改修 + ミドルウェア稼働修正

### レスポンシブデザイン改修
- **`#app` が全幅にならない問題**: Viteスターター残骸の `body { display:flex; place-items:center }` により、アプリ全体がコンテンツ幅で中央寄せされていた。`#app { width:100% }` に修正（全画面に影響する根本修正）
- **Dashboardナビゲーション**: モバイル(<768px)ではボタンをアイコンのみに圧縮する `.nav-compact` ユーティリティを導入し、折返し可能に（ラベルはDOMに残るためアクセシビリティ維持）
- **固定幅ダイアログ**: 400〜500px固定だった全ダイアログ（Tasks編集/一括追加、Projects作成、Teams作成、Calendar作成/編集、TeamDetailリネーム）に `:breakpoints="{'640px':'92vw'}"` を追加し、スマートフォンでの画面外はみ出しを解消
- **行の折返し**: タスクカードの操作ボタン、一括操作バー、チームメンバー行の操作群に flex-wrap を追加
- **Priority選択ボタンのラベル欠落**（既存バグ）: SelectButton に関数 optionLabel を渡しておりラベルが描画されなかった → オブジェクト配列 + optionLabel/optionValue に修正（Tasks / Calendar）
- **検証**: 375px / 1280px のスクリーンショットで Login・Dashboard・Tasks・編集ダイアログを確認。vue-tsc ビルド / vitest / Playwright E2E 5/5 全パス

### ミドルウェア稼働修正（docker-compose が成立していなかった）
- **frontendイメージのビルド失敗**: コンテナは `npm i -g pnpm` で pnpm 11 が入り、pnpm 10+ のビルドスクリプト承認ゲートで `pnpm install` がエラー終了 → `packageManager: pnpm@10.33.2` + corepack に切替え、`onlyBuiltDependencies`（esbuild）/`ignoredBuiltDependencies` を package.json と pnpm-workspace.yaml の両方に定義。Dockerfile は `--frozen-lockfile` 化
- **コンテナ内viteプロキシの宛先誤り**: フロントコンテナの `/api` プロキシが自身の `localhost:8888` を向いており、Docker構成ではバックエンドへ到達不能だった → compose に `VITE_BACKEND_URL=http://backend:8888` を追加（コンテナ経由のログイン〜E2Eで動作確認済み）
- **DBホストポートの衝突**: 他プロジェクトが 5432 を占有していると `docker-compose up` が失敗 → `"${DB_PORT:-5432}:5432"` にパラメータ化（`DB_PORT=15432 docker-compose up` で回避可能。コンテナ間通信には影響なし）
- **JWTミドルウェア**: `/schemas/*`（OpenAPIのJSON Schema参照。エラーレスポンスの `$schema` や /docs が参照）が401になっていた → 公開パスに追加
- **稼働確認**: db (pg_isready) / backend (openapi・docs・schemas 200、未認証401、CORSプリフライト200+Allow-Origin) / frontend (5173応答、プロキシ経由ログイン成功) すべて正常

---

## 2026-06-13 ガントチャート画面の検証と改修

### 発覚した問題（Playwrightで実操作して検証）
1. **チャート全体が真っ黒に表示される**: frappe-gantt のスタイルシートが一度もimportされていなかった（コード内に「要調整」コメントのまま残置）。SVG要素がデフォルトの黒塗りで描画され、グリッド・ラベル・バーが一切視認できない状態だった
2. **バーのドラッグ操作が機能しない**: ドラッグ中に `Cannot read properties of undefined (reading 'classList')` 例外が連発し、`on_date_change` が発火せず期日が保存されなかった。原因は1と同根 — スタイル未適用でヘッダーラベルのレイアウト幅が壊れ、frappe-gantt 内部のスクロールハンドラの日付ラベル探索が失敗していた

### 対応内容
- スタイルシートを読み込み（package の `exports` がdeep importを塞ぎ、postcss-import は `style` 条件を解決しないため、`@import "../../node_modules/frappe-gantt/dist/frappe-gantt.css"` でSFCのstyleブロックから直接参照）
- Gantt インスタンスを reactive ref に格納しない構造へ変更（Vueのプロキシがライブラリ内部のDOM参照を包むのを回避）
- `view_mode_select: true` を追加し Day / Week / Month をUIから切替可能に。v0.x時代の無効オプション（`custom_popup_html` 等）を整理
- **対象ファイル**: `frontend/src/views/Gantt.vue`

### 検証結果（実ブラウザ操作）
- 表示: 月/日ヘッダ・本日ハイライト・進捗付きバー・ラベルが正しく描画
- 操作: バーを+3日ドラッグ → `PUT /tasks/{id}` 200 → トースト表示 → APIで期日+3日を確認 → リロード後もバー位置が維持。コンソールエラー0件
- ビューモード切替（Day/Week/Month）、バークリックのポップアップ表示も正常
- 回帰: `pnpm build` / vitest 6/6 / Playwright E2E 5/5 全パス
