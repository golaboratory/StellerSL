# StellerSL

## 1. 開発の目的

PC、タブレット、スマートフォンなど、あらゆるデバイスから簡易に入力・閲覧できるよう、**レスポンシブデザインを前提**としたTODO管理アプリを開発する。
複数のプロジェクトを管理可能で、各プロジェクトには複数のタスクが含まれる。
各タスクは、一括で登録やステータス変更が可能とする。

また、自身のタスクの一覧では、タスクの総数や、未完了のタスク数、完了のタスク数が表示される。
一日のタスクの消化数や、登録数でグラフを表示する。
モチベーション維持のため、達成率に応じてバッジの獲得や、キャラクターの成長などの要素も含む。

個人のタスク一覧では、GTDの考え方を取り入れ、簡易にステータスの変更を行えるようにする。

マルチテナント方式を採用し、アクセスされるドメイン名によってテナントを識別する。
データアーキテクチャとしては各テーブルに `tenant_id` (UUID) を持たせることで、DB内でテナント間のデータを**論理分離**している。

### 1-1. 機能一覧 (実装済み)

- **ユーザー管理**: ユーザー登録 (Signup)、ログイン (JWT発行)、ログアウト
- **プロジェクト管理**: プロジェクト一覧、登録、更新、削除、ユーザーアサイン
- **タスク管理 (GTD)**: タスク一覧 (簡易ステータス変更)、詳細、登録、更新、削除
- **一括操作**: タスク一括登録 (行区切り入力)、一括ステータス変更、一括削除
- **ダッシュボード**: タスク統計表示、日別アクティビティグラフ (直近7日間)
- **チーム機能**: チーム作成、チームメンバー管理 (招待リンク/サインアップ連携)
- **ゲーミフィケーション**: キャラクター成長システム (10 EXP/task, 100 EXP/level)、バッジ獲得・表示
- **高度な可視化**: 
  - ガントチャート表示 (frappe-gantt)
  - カレンダー表示 (v-calendar)

## 2. 技術スタック

- **フロントエンド**
  - Vue 3 (Composition API / `<script setup>`)
  - TypeScript
  - Tailwind CSS v4 (Vite plugin)
  - PrimeVue 4 (Aura theme / @primeuix/themes)
  - OpenAPI Generator (Axios client)
  - Pinia (状態管理 / localStorage 永続化)
  - pnpm
- **バックエンド**
  - Go
  - Huma v2 (OpenAPI 3.1 自動生成)
  - PostgreSQL (lib/pq)
  - JWT (Authentication)
  - Bcrypt (Password Hashing)
- **インフラ・開発ツール**
  - Docker / Docker Compose
  - Air (Go hot-reloading)
  - PowerShell / Bash (API生成スクリプト)

## 3. 開発の進め方 (スキーマ駆動)

### 3-1. DBアクセス層 (sqlc)
- `backend/internal/db/**` 配下のファイルは `sqlc` によって自動生成されるため、**直接編集してはいけない**。
- クエリを追加・変更する場合は、`backend/sql/query.sql` を編集し、`sqlc generate` (または `docker-compose exec backend sqlc generate`) を実行すること。
- RLS (Row Level Security) を考慮し、セッション変数 `app.current_tenant_id` を利用するクエリ設計を行う。

### 3-2. APIクライアント生成 (OpenAPI)
- `frontend/src/api/**` 配下のファイルは `OpenAPI Generator` によって自動生成されるため、**直接編集してはいけない**。
- API定義が変更された場合は、バックエンドを起動した状態で `.\generate-api.ps1` (Windows) または `./generate-api.sh` (Unix) を実行すること。
- スクリプト内部では、起動中のバックエンドから `openapi.json` を取得し、`openapi-generator-cli` を用いてクライアントコードを再生成する。

### 3-3. API更新サイクル
1. バックエンドで `huma` を用いてエンドポイントを実装/更新。
2. 上記の生成スクリプトを実行して API クライアントを更新。
3. `frontend/src/api` に生成された型安全なクラスを利用して UI を実装。


## 4. ディレクトリ構造

```text
.
├── backend/            # Go (huma) バックエンド
│   ├── cmd/            # エントリポイント (main.go)
│   ├── internal/
│   │   ├── api/        # Huma API ハンドラー・ルート定義
│   │   └── db/         # データベースアクセス層 (CRUDロジック)
│   └── sql/            # DDLおよびクエリ定義
├── frontend/           # Vue 3 フロントエンド
│   ├── src/
│   │   ├── api/        # 生成されたAPIクライアント
│   │   ├── components/ # 共通UI
│   │   ├── stores/     # Pinia (Auth等)
│   │   └── views/      # 画面コンポーネント (Dashboard, Tasks, etc)
├── _initdb/            # DB初期投入スクリプト (01_schema.sql, etc)
└── docker-compose.yml  # 実行環境定義
```

## 5. テストの方針

- **フロントエンド**: Vitest を使用。Pinia ストアのロジックやコンポーネントのテスト。
- **バックエンド**: 標準の `testing` パッケージによるロジック検証。

## 6. 特記事項 (実装仕様)

- **テナント識別**: `Host` ヘッダーからドメインを取得し、`GetTenantByDomain` クエリで `tenant_id` を特定する。
- **EXP計算**: タスク完了時に 10 EXP 付与。100 EXP ごとにレベルアップ (`level = (exp / 100) + 1`)。
- **一括登録**: テキストエリアによる1行1タスクの簡易入力を採用。
