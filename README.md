# Yelp Sample API v2

Yelp風のレストラン検索・レビューAPIのサンプルプロジェクトです。Go言語とGinフレームワークを使用し、PostgreSQLデータベースと連携しています。

## 機能

- **ビジネス検索**: 位置情報やカテゴリでレストランを検索
- **ビジネス詳細**: 個別のレストラン情報を取得
- **レビュー機能**: レストランのレビューを投稿・取得
- **RESTful API**: 標準的なHTTP APIエンドポイント

## 技術スタック

- **Backend**: Go 1.23
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Containerization**: Docker & Docker Compose

## API エンドポイント

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | サーバーステータス確認 |
| GET | `/health` | ヘルスチェック |
| GET | `/businesses` | ビジネス検索 |
| GET | `/businesses/:id` | ビジネス詳細取得 |
| GET | `/businesses/:id/reviews` | ビジネスのレビュー取得 |
| POST | `/businesses/:id/reviews` | レビュー投稿 |

## セットアップ

### 必要条件

- Docker
- Docker Compose

### 起動方法

1. リポジトリをクローン
```bash
git clone <repository-url>
cd yelp_sample_v2
```

2. Docker Composeでサービスを起動
```bash
docker-compose up -d
```

3. データベースの初期化（任意）
```bash
# サンプルデータを投入
docker exec -i yelp_postgres psql -U postgres -d yelp_sample < sample_data.sql
```

4. サーバーの動作確認
```bash
curl http://localhost:8080/health
```

## 環境変数

アプリケーションは以下の環境変数を使用します：

- `DB_HOST`: データベースホスト（デフォルト: postgres）
- `DB_USER`: データベースユーザー（デフォルト: postgres）
- `DB_PASSWORD`: データベースパスワード（デフォルト: postgres）
- `DB_NAME`: データベース名（デフォルト: yelp_sample）
- `DB_PORT`: データベースポート（デフォルト: 5432）
- `PORT`: アプリケーションポート（デフォルト: 8080）

## データベーススキーマ

### Users テーブル
- ユーザー情報（ID、名前、メール、作成日時、更新日時）

### Businesses テーブル
- ビジネス情報（ID、名前、カテゴリ、位置情報、住所、電話番号、ウェブサイト、説明、作成日時、更新日時）

### Reviews テーブル
- レビュー情報（ID、ビジネスID、ユーザーID、評価、テキスト、作成日時、更新日時）

## 開発

### ローカル開発環境での実行

```bash
# 依存関係のインストール
go mod download

# データベースの起動
docker-compose up -d postgres

# アプリケーションの起動
go run main.go
```

### テストの実行

```bash
go test ./...
```

## サンプルデータ

`sample_data.sql`には以下のサンプルデータが含まれています：

- 5名のユーザー
- 8件のビジネス（ラーメン店、カフェ、寿司店、イタリアン、居酒屋、焼肉店、ベーカリー、中華料理）
- 20件のレビュー

## API使用例

### ビジネス検索
```bash
curl "http://localhost:8080/businesses?category=ラーメン&limit=10"
```

### ビジネス詳細取得
```bash
curl "http://localhost:8080/businesses/1"
```

### レビュー取得
```bash
curl "http://localhost:8080/businesses/1/reviews"
```

### レビュー投稿
```bash
curl -X POST "http://localhost:8080/businesses/1/reviews" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "rating": 5, "text": "とても美味しかったです！"}'
```

## プロジェクト構成

```
yelp_sample_v2/
├── main.go                 # アプリケーションエントリーポイント
├── go.mod                  # Go モジュール設定
├── go.sum                  # 依存関係チェックサム
├── Dockerfile              # Docker設定
├── docker-compose.yml      # Docker Compose設定
├── sample_data.sql         # サンプルデータ
├── db_queries.sql          # データベースクエリ
├── database/               # データベース設定
│   └── database.go
├── models/                 # データモデル
│   ├── business.go
│   ├── review.go
│   └── user.go
└── handlers/               # HTTPハンドラー
    └── business.go
```