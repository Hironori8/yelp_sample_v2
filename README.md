# Yelp Sample API v2 - Microservices Architecture

Yelp風のレストラン検索・レビューAPIのマイクロサービス版です。APIゲートウェイを通じて、ビジネスサービスとレビューサービスが独立して動作します。

## アーキテクチャ

本プロジェクトはマイクロサービスアーキテクチャを採用しており、以下のサービスで構成されています：

- **APIゲートウェイ** (Port 8080): リクエストを各サービスにルーティング
- **ビジネスサービス** (Port 8081): 店舗情報の管理
- **レビューサービス** (Port 8082): レビューの管理
- **PostgreSQL データベース**: 全サービスで共有される単一データベース

## 機能

- **ビジネス検索**: 位置情報やカテゴリでレストランを検索
- **ビジネス詳細**: 個別のレストラン情報を取得
- **レビュー機能**: レストランのレビューを投稿・取得
- **RESTful API**: 標準的なHTTP APIエンドポイント
- **マイクロサービス**: 独立してスケール可能なサービス分離

## 技術スタック

- **Backend**: Go 1.23
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Architecture**: Microservices with API Gateway
- **Containerization**: Docker & Docker Compose

## API エンドポイント

すべてのリクエストはAPIゲートウェイ (http://localhost:8080) 経由で行います。

| Method | Endpoint | Service | Description |
|--------|----------|---------|-------------|
| GET | `/` | Gateway | APIゲートウェイステータス確認 |
| GET | `/health` | Gateway | ヘルスチェック |
| GET | `/businesses` | Business | ビジネス検索 |
| GET | `/businesses/:id` | Business | ビジネス詳細取得 |
| GET | `/businesses/:id/reviews` | Review | ビジネスのレビュー取得 |
| POST | `/businesses/:id/reviews` | Review | レビュー投稿 |
| GET | `/reviews` | Review | 全レビュー取得 |
| GET | `/reviews/:id` | Review | 個別レビュー取得 |

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

2. Docker Composeで全サービスを起動
```bash
docker-compose up --build -d
```

3. サービスの起動確認
```bash
docker-compose ps
```

4. データベースの初期化確認
```bash
# サンプルデータが既に投入されていることを確認
docker exec yelp_postgres psql -U postgres -d yelp_sample -c "\dt"
```

5. APIゲートウェイの動作確認
```bash
curl http://localhost:8080/health
```

## サービス詳細

### APIゲートウェイ
- **Port**: 8080
- **役割**: リクエストのルーティング、負荷分散
- **エンドポイント**: 全APIのエントリーポイント

### ビジネスサービス
- **Port**: 8081（内部通信用）
- **役割**: 店舗情報の管理
- **データベース**: PostgreSQL（共有）

### レビューサービス
- **Port**: 8082（内部通信用）
- **役割**: レビューデータの管理
- **データベース**: PostgreSQL（共有）

## 環境変数

各サービスは以下の環境変数を使用します：

### データベース設定
- `DB_HOST`: データベースホスト（デフォルト: postgres）
- `DB_USER`: データベースユーザー（デフォルト: postgres）
- `DB_PASSWORD`: データベースパスワード（デフォルト: postgres）
- `DB_NAME`: データベース名（デフォルト: yelp_sample）
- `DB_PORT`: データベースポート（デフォルト: 5432）

### サービス固有設定
- `PORT`: 各サービスのポート番号
  - Gateway: 8080
  - Business: 8081
  - Review: 8082

## データベーススキーマ

### Users テーブル
- ユーザー情報（ID、名前、メール、作成日時、更新日時）

### Businesses テーブル
- ビジネス情報（ID、名前、カテゴリ、位置情報、住所、電話番号、ウェブサイト、説明、評価、レビュー数、作成日時、更新日時）

### Reviews テーブル
- レビュー情報（ID、ビジネスID、ユーザーID、評価、テキスト、作成日時、更新日時）

## 開発

### ローカル開発環境での実行

個別サービスの開発時：

```bash
# データベースのみ起動
docker-compose up -d postgres

# ビジネスサービスの起動
cd services/business
go mod download
go run main.go

# レビューサービスの起動（別ターミナル）
cd services/review
go mod download
go run main.go

# APIゲートウェイの起動（別ターミナル）
cd services/gateway
go mod download
go run main.go
```

### ログの確認

各サービスのログを確認：

```bash
# 全サービスのログ
docker-compose logs

# 特定サービスのログ
docker logs yelp_gateway
docker logs yelp_business_service
docker logs yelp_review_service
docker logs yelp_postgres
```

## サンプルデータ

`sample_data.sql`には以下のサンプルデータが含まれています：

- 5名のユーザー
- 8件のビジネス（ラーメン店、カフェ、寿司店、イタリアン、居酒屋、焼肉店、ベーカリー、中華料理）
- 21件のレビュー

### 現在のデータ状況
```bash
# データ件数確認
docker exec yelp_postgres psql -U postgres -d yelp_sample -c "SELECT COUNT(*) FROM users;"
docker exec yelp_postgres psql -U postgres -d yelp_sample -c "SELECT COUNT(*) FROM businesses;"
docker exec yelp_postgres psql -U postgres -d yelp_sample -c "SELECT COUNT(*) FROM reviews;"
```

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
  -d '{"rating": 5, "text": "とても美味しかったです！"}'
```

### 全レビュー取得
```bash
curl "http://localhost:8080/reviews"
```

## プロジェクト構成

```
yelp_sample_v2/
├── docker-compose.yml           # Docker Compose設定
├── sample_data.sql              # サンプルデータ
├── db_queries.sql               # データベースクエリ
├── services/                    # マイクロサービス
│   ├── gateway/                 # APIゲートウェイ
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── Dockerfile
│   ├── business/                # ビジネスサービス
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Dockerfile
│   │   ├── database/            # データベース設定
│   │   │   └── database.go
│   │   ├── models/              # データモデル
│   │   │   ├── business.go
│   │   │   ├── review.go
│   │   │   └── user.go
│   │   └── handlers/            # HTTPハンドラー
│   │       └── business.go
│   └── review/                  # レビューサービス
│       ├── main.go
│       ├── go.mod
│       ├── go.sum
│       ├── Dockerfile
│       ├── database/            # データベース設定
│       │   └── database.go
│       ├── models/              # データモデル
│       │   ├── business.go
│       │   ├── review.go
│       │   └── user.go
│       └── handlers/            # HTTPハンドラー
│           └── review.go
└── <legacy files in root>/      # 旧モノリシック版（参考用）
    ├── main.go
    ├── database/
    ├── models/
    └── handlers/
```

## トラブルシューティング

### ポート競合エラー
```bash
# 既存のコンテナを停止
docker-compose down
docker container prune

# 再起動
docker-compose up --build -d
```

### データベース接続エラー
```bash
# PostgreSQLコンテナの状態確認
docker logs yelp_postgres

# データベースの健康状態確認
docker-compose ps
```

### サービス間通信エラー
```bash
# ネットワーク設定確認
docker network ls
docker network inspect yelp_sample_v2_default
```

## スケーリング

個別サービスのスケーリング例：

```bash
# レビューサービスを3インスタンスにスケール
docker-compose up --scale review-service=3 -d

# ビジネスサービスを2インスタンスにスケール
docker-compose up --scale business-service=2 -d
```

## ステータス

✅ **完了**
- マイクロサービスアーキテクチャの実装
- Docker Compose環境の構築
- API Gateway の実装
- Business Service の実装
- Review Service の実装
- PostgreSQL データベース設定
- サンプルデータの投入
- 全サービスの動作確認完了

## 今後の拡張予定

- [ ] サービス間認証の実装
- [ ] API Gateway での負荷分散
- [ ] サービスディスカバリーの導入
- [ ] 分散トレーシングの実装
- [ ] メトリクス監視の追加
- [ ] フロントエンド（Next.js）の追加