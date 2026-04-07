# TODO.md

未実装タスク一覧。完了済みタスクは [COMPLETE.md](COMPLETE.md) を参照。

---

## Phase 5: タスク機能の拡充

### 5-1. タスクのフィルタリング機能
- **現状**: タスク一覧はフロントエンドでの検索のみ（タイトル/説明の部分一致）。API側のフィルタがない
- **対応**: `GET /tasks` にクエリパラメータ (`status`, `priority`, `project_id`, `due_date_from`, `due_date_to`) を追加
- **対象ファイル**: `backend/internal/api/task/handler.go`, `backend/sql/query.sql`

### 5-2. タスク担当者の変更機能
- **現状**: `assigned_to` はタスク作成時に設定されるが、後から変更するUIやエンドポイントがない
- **対応**: `PUT /tasks/{id}` で `assigned_to` を更新可能にする（既存エンドポイントの拡張）

---

## Phase 6: フロントエンド品質改善

### 6-1. フォームバリデーションの強化
- **現状**: HTML5の `required` 属性のみ。メールフォーマット検証、パスワード強度チェック、カスタムエラー表示がない
- **対応**: PrimeVueのフォームバリデーション機能を活用し、各フォームにバリデーションを追加
- **対象ファイル**: `Login.vue`, `Signup.vue`, `Tasks.vue`（タスク編集ダイアログ）, `Projects.vue`

### 6-3. ガントチャートの編集機能実装
- **現状**: ガントチャート上でドラッグしても `console.log` のみでAPIへの反映がない
- **対応**: 日付変更時にAPIの `PUT /tasks/{id}` を呼び出してdue_dateを更新する
- **対象ファイル**: `frontend/src/views/Gantt.vue`

### 6-4. カレンダービューからのタスク編集
- **現状**: カレンダーはタスクの閲覧のみで、日付をクリックしてもタスクの編集や作成ができない
- **対応**: 日付クリック時のタスク作成/編集ダイアログを追加
- **対象ファイル**: `frontend/src/views/CalendarView.vue`

### 6-6. チームメンバー追加のUX改善
- **現状**: チームメンバー追加はUUID手入力が必要（`TeamDetail.vue`）
- **対応**: プロジェクトメンバー追加と同様のユーザー検索UIに変更する
- **対象ファイル**: `frontend/src/views/TeamDetail.vue`

---

## Phase 7: テストカバレッジの拡充

### 7-1. バックエンド: サービス層のユニットテスト追加
- **現状**: バックエンドのテストはルーティング検証（OpenAPIスペック確認）がほとんどで、実際のビジネスロジックテストがない
- **対応**: 各domain (task, project, team) のservice層にモックDBを使ったユニットテストを追加

### 7-2. バックエンド: バッジエンジンのユニットテスト
- **現状**: バッジ判定ロジックのテストが一切ない
- **対応**: 全バッジ種別に対する判定テストを追加
- **対象ファイル**: `backend/internal/api/gamification/engine_test.go`（新規作成）

### 7-3. フロントエンド: Vitest単体テストの拡充
- **現状**: auth store と router のみ（計6テスト）
- **対応**: toast store, 主要コンポーネントのレンダリングテストを追加

### 7-4. E2Eテスト: チーム機能
- **現状**: チーム関連のE2Eテストが一切ない
- **対応**: チーム作成、メンバー追加/削除、招待リンクのE2Eテストを追加

### 7-5. バックエンド: マルチテナント分離テスト
- **現状**: テナント分離がテストで検証されていない
- **対応**: 異なるテナントIDでのデータアクセスが正しく分離されることを検証するテストを追加

---

## Phase 8: インフラ・運用改善

### 8-2. 構造化ログの導入
- **現状**: `fmt.Println` / `fmt.Printf` のみ
- **対応**: `slog` (Go標準) またはサードパーティのロガーを導入し、JSON形式の構造化ログに切り替え

### 8-3. 本番環境のデプロイ構成
- **決定事項**: Q4=B — VPS + Docker Compose + Let's Encrypt
- **対応**: `docker-compose.prod.yml` を作成。Nginx リバースプロキシ + Let's Encrypt + prod Dockerfile。`.env.example` で環境変数ドキュメント

---

## フロントエンド: API再生成後に実装（Sprint 7 残作業）

> 以下はバックエンドを起動し `./generate-api.sh` で OpenAPI クライアントを再生成した後に着手可能

### 7-FE-1. チームメンバー検索UI
- `TeamDetail.vue`: UUID手入力 → ユーザー検索オートコンプリートに変更

### 7-FE-2. チームロール管理UI
- `TeamDetail.vue`: ロール表示・変更ドロップダウン、チーム編集/削除ボタン

### 7-FE-3. 通知ベル + ドロワー
- `NotificationBell.vue`, `NotificationDrawer.vue` (新規), `stores/notification.ts` (新規)
- 30秒ポーリングで未読カウント更新、ドロワーで一覧表示

### 7-FE-4. プロフィールにパスワード変更UI
- `Profile.vue`: 現在/新規/確認パスワード入力 + PUT /auth/password 呼び出し

### 7-FE-5. キャラクター進化表示
- `Dashboard.vue`: character_type に応じたアイコン・名称表示 (egg/chick/chicken/phoenix)

### 7-FE-6. ガントチャート編集 → 6-3 と同一

### 7-FE-7. カレンダーからタスク作成 → 6-4 と同一
