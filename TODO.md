# StellerSL Development TODO

## 1. 基盤構築・PoC (Phase 1) - [COMPLETE]
- [x] Docker Compose によるローカル開発環境の構築
- [x] DBスキーマ設計および初期投入スクリプト (_initdb)
- [x] Go + Huma による API 基盤のセットアップ
- [x] Vue 3 + TypeScript + Tailwind CSS v4 + PrimeVue 4 の導入
- [x] OpenAPI Generator による型付き API クライアント生成自動化 (.ps1 / .sh)

## 2. 認証・認可 (Authentication) - [COMPLETE]
- [x] ユーザーログイン (JWT発行)
- [x] ユーザー登録 (サインアップ / Bcrypt)
- [x] ログアウト (クライアント側トークン破棄)
- [x] テナント識別ミドルウェアの完全実装 (Hostヘッダーベース)
- [x] PostgreSQL Row Level Security (RLS) による厳密なテナント分離の実装
- [x] JWT 有効期限切れの自動ハンドリング (401 Axios Interceptor)

## 3. プロジェクト管理 (Project Management)
- [x] プロジェクト一覧表示
- [x] プロジェクト登録 (API & UI)
- [x] プロジェクト更新 (API)
- [x] プロジェクト削除 (API & UI)
- [x] プロジェクト詳細画面の実装 (Task絞り込み連携)
- [x] プロジェクトへのユーザーアサイン管理 UI (API & UI 実装済み)


## 4. タスク管理 (Task Management)
- [x] タスク一覧表示 (GTD要素を含む簡易ステータス変更)
- [x] タスク登録 (単一)
- [x] タスクステータス更新 (PATCH /tasks/{id}/status)
- [x] タスク削除 (単一)
- [x] タスク一括登録 (行区切りテキスト入力)
- [x] タスク一括ステータス更新
- [x] タスク一括削除
- [x] タスク編集モーダルの実装 (タイトル・説明・期限の詳細編集)
- [x] 優先度(Priority)のUI反映 (アイコン表示)

## 5. ダッシュボード・統計 (Dashboard & Analytics) - [COMPLETE]
- [x] タスク統計 (Total/Pending/Completed)
- [x] 日別アクティビティグラフ (7日間)
- [x] ゲーミフィケーション情報の統合表示

## 6. チーム機能 (Team Features) - [COMPLETE]
- [x] チーム作成・管理
- [x] チームメンバー管理 (API)
- [x] チーム詳細 UI (メンバー一覧表示、追加フォーム)
- [x] チーム招待リンク生成機能 (Signupクエリパラメータ連携)

## 7. ゲーミフィケーション (Gamification) - [COMPLETE]
- [x] レベル・経験値システム (10 EXP/task)
- [x] バッジ獲得・一覧表示 (API & UI)
- [x] キャラクター成長に応じたアバター変化 (レベルに応じたアイコン・色)
- [x] 特殊バッジ獲得通知 (Dashboard Toast連携)

## 8. 高度な表示 (Advanced Views)
- [x] ガントチャート基本表示 (frappe-gantt)
- [x] カレンダー基本表示 (v-calendar)
- [x] ガントチャートでの日付変更イベント検知
- [x] カレンダーでの日別タスク表示連携

## 9. テスト・品質 (Quality Assurance)
- [x] バックエンド: DBアクセス層・API層の構築
- [x] フロントエンド: Pinia ストアの単体テスト
- [x] バックエンド: ユニットテスト (api ヘルパー関数のテスト、ハンドラーのルーティングテスト)
- [x] フロントエンド: 画面遷移テスト (Vue Router / Navigation Guards)
- [x] 継続的インテグレーション (GitHub Actions)

## 10. 今後の課題 (Next Steps) - [COMPLETE]
- [x] 本番環境用 Dockerfile の最適化 (マルチステージビルド)
- [x] データベースマイグレーションのディレクトリ構造化
- [x] フロントエンドのエラーハンドリング共通化 (Global Toast Interceptor)
- [x] ダークモード/ライトモードの切り替え UI

## 11. 機能拡張・品質向上 (Phase 2 Enhancement)
- [x] バックエンド: プロジェクト詳細、更新、削除の API ハンドラー実装 (huma)
- [x] バックエンド: プロジェクト別タスク取得 API (`GET /projects/{id}/tasks`)
- [x] バックエンド: プロジェクト・ユーザーアサイン API の実装 (`POST /projects/{id}/users`, `DELETE /projects/{id}/users/{user_id}`)
- [x] バックエンド: チームメンバー一覧取得 API (`GET /teams/{id}/members`)
- [x] フロントエンド: `ProjectDetail.vue` のリファクタリング (全件取得フィルタリングから個別取得へ)
- [x] フロントエンド: プロジェクトメンバー管理 UI の完全実装
- [x] フロントエンド: チームメンバー管理 UI の補完 (一覧表示、メンバー削除)
- [x] フロントエンド: ユーザープロフィール編集機能 (名前、アバター URL)
- [x] フロントエンド: タスクのグローバル検索・フィルタリング機能
- [x] バックエンド: タスク完了時のバッジ自動付与ロジックの強化
- [x] バックエンド: アクティビティログの取得・表示 API (Dashboard詳細)

## 12. 機能拡張・品質向上 (Phase 3 Enhancement)
- [x] バックエンド: タスク詳細更新 API の実装 (`PUT /tasks/{id}`)
- [x] フロントエンド: タスク編集モーダルの全フィールド保存対応 (API 連携)
- [x] フロントエンド: ダッシュボードでの最近のアクティビティ一覧表示
- [x] バックエンド: リスト取得系 API へのページネーション導入 (`projects`, `tasks`, `teams`)
- [x] バックエンド: バッジエンジン（条件判定ロジック）の独立化と強化
- [ ] フロントエンド: バリデーションエラー (400 Bad Request) の詳細表示
- [x] フロントエンド: プロジェクトメンバー追加時のユーザー検索 (Email/Name検索)
- [ ] フロントエンド: アバター画像のファイルアップロード対応

## 13. 仕様確認 (Specification Clarification)
- [x] ページネーションのデフォルト件数および最大件数の決定　→ デフォルトを50件、最大を200件とする
- [x] バッジ獲得条件の具体的なリスト化 (例: 10タスク完了、3日間連続ログイン等)　→ /_memo/バッジの仕様.md に記載
- [x] チーム招待リンクの有効期限設定の要否　→ 1週間とする
- [x] 削除済みタスクの論理削除/物理削除の方針確認 → 論理削除とする

