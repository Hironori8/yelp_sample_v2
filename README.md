# Yelp Sample API v2 - Microservices Architecture

Yelp風のレストラン検索・レビューAPIのマイクロサービス版です。APIゲートウェイを通じて、認証、ビジネス、レビュー、ログサービスが独立して動作します。

## アーキテクチャ

本プロジェクトはマイクロサービスアーキテクチャを採用しており、以下のサービスで構成されています：

- **APIゲートウェイ** (Port 8080): 認証とリクエストルーティング
- **認証サービス**: ユーザー認証・JWT発行
- **ビジネスサービス**: 店舗情報の管理
- **レビューサービス**: レビューの管理
- **ログサービス**: ユーザー行動ログの記録
- **PostgreSQL データベース**: ユーザー・ビジネス・レビューデータ
- **Cassandra データベース**: ログデータ（時系列データ）

### ネットワーク構成

- **Frontend Network**: APIゲートウェイのみ外部公開 (Port 8080)
- **Backend Network**: 内部サービス間通信専用（外部アクセス不可）

## 機能

- **ユーザー認証**: JWT認証による安全な認証システム
- **ビジネス検索**: 位置情報やカテゴリでレストランを検索
- **ビジネス詳細**: 個別のレストラン情報を取得
- **レビュー機能**: 認証ユーザーがレビューを投稿・取得
- **ユーザー行動ログ**: レビュー閲覧履歴の記録・分析
- **RESTful API**: 標準的なHTTP APIエンドポイント
- **マイクロサービス**: 独立してスケール可能なサービス分離

## 技術スタック

- **Backend**: Go 1.23
- **Web Framework**: Gin
- **ORM**: GORM (PostgreSQL), ScyllaDB/Cassandra (ログデータ)
- **Authentication**: JWT (JSON Web Token)
- **Database**: 
  - PostgreSQL 15 (トランザクショナルデータ)
  - Cassandra 4.0 (ログ・時系列データ)
- **Architecture**: Microservices with API Gateway
- **Containerization**: Docker & Docker Compose

## API エンドポイント

すべてのリクエストはAPIゲートウェイ (http://localhost:8080) 経由で行います。

### パブリックエンドポイント（認証不要）

| Method | Endpoint | Service | Description |
|--------|----------|---------|-------------|
| GET | `/` | Gateway | APIゲートウェイステータス確認 |
| GET | `/health` | Gateway | ヘルスチェック |
| POST | `/auth/register` | Auth | ユーザー登録 |
| POST | `/auth/login` | Auth | ログイン（JWT取得） |
| POST | `/auth/logout` | Auth | ログアウト |
| GET | `/businesses` | Business | ビジネス検索 |
| GET | `/businesses/:id` | Business | ビジネス詳細取得 |
| GET | `/businesses/:id/reviews` | Review | ビジネスのレビュー取得 |

### 保護されたエンドポイント（JWT認証必要）

| Method | Endpoint | Service | Description |
|--------|----------|---------|-------------|
| GET | `/auth/me` | Auth | 現在ユーザー情報取得 |
| POST | `/businesses/:id/reviews` | Review | レビュー投稿 |
| GET | `/reviews` | Review | 全レビュー取得 |
| GET | `/reviews/:id` | Review | 個別レビュー取得 |
| GET | `/logs/user/:user_id/history` | Logging | ユーザー閲覧履歴 |
| GET | `/logs/business/:business_id/stats` | Logging | ビジネス統計 |

## セットアップ

### 必要条件

- Docker
- Docker Compose

### 起動方法

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
# PostgreSQLのサンプルデータ確認
docker exec yelp_postgres psql -U postgres -d yelp_sample -c "\dt"

# Cassandraのキースペース確認
docker exec yelp_cassandra cqlsh -e "DESCRIBE KEYSPACES;"
```

5. APIゲートウェイの動作確認
```bash
curl http://localhost:8080/health
```

## サービス詳細

### APIゲートウェイ
- **Port**: 8080
- **役割**: JWT認証、リクエストルーティング、セキュリティ制御
- **認証方式**: ヘッダーベースのユーザー情報伝達
- **エンドポイント**: 全APIのエントリーポイント

### 認証サービス
- **役割**: ユーザー登録・ログイン・JWT発行
- **データベース**: PostgreSQL（ユーザー情報）
- **セキュリティ**: bcrypt パスワードハッシュ化

### ビジネスサービス
- **役割**: 店舗情報の管理
- **データベース**: PostgreSQL（共有）

### レビューサービス
- **役割**: レビューデータの管理、ユーザー行動ログ送信
- **データベース**: PostgreSQL（共有）

### ログサービス
- **役割**: ユーザー行動ログの記録・分析
- **データベース**: Cassandra（時系列データ）

## 環境変数

各サービスは以下の環境変数を使用します：

### データベース設定（PostgreSQL）
- `DB_HOST`: データベースホスト（デフォルト: postgres）
- `DB_USER`: データベースユーザー（デフォルト: postgres）
- `DB_PASSWORD`: データベースパスワード（デフォルト: postgres）
- `DB_NAME`: データベース名（デフォルト: yelp_sample）
- `DB_PORT`: データベースポート（デフォルト: 5432）

### 認証設定
- `JWT_SECRET`: JWT署名シークレット
- `JWT_EXPIRES_IN`: JWT有効期限（デフォルト: 24h）

### サービス固有設定
- `PORT`: APIゲートウェイのポート番号（8080）

### ログサービス設定
- `CASSANDRA_HOSTS`: Cassandraホスト（デフォルト: cassandra:9042）

## データベーススキーマ

### PostgreSQL（トランザクショナルデータ）

#### Users テーブル
- ユーザー情報（ID、名前、メール、パスワード、作成日時、更新日時）

#### Businesses テーブル
- ビジネス情報（ID、名前、カテゴリ、位置情報、住所、電話番号、ウェブサイト、説明、評価、レビュー数、作成日時、更新日時）

#### Reviews テーブル
- レビュー情報（ID、ビジネスID、ユーザーID、評価、テキスト、作成日時、更新日時）

### Cassandra（ログデータ）

#### review_view_logs テーブル
- レビュー閲覧ログ（ユーザーID、ビジネスID、レビューID、閲覧日時、IPアドレス、ユーザーエージェント）

## 認証フロー

### 1. ユーザー登録
```bash
curl -X POST "http://localhost:8080/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"山田太郎","email":"yamada@example.com","password":"password123"}'
```

### 2. ログイン（JWTトークン取得）
```bash
TOKEN=$(curl -s -X POST "http://localhost:8080/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"yamada@example.com","password":"password123"}' | jq -r '.token')
```

### 3. 認証が必要なAPIの呼び出し
```bash
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/reviews"
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

### レビュー取得（パブリック）
```bash
curl "http://localhost:8080/businesses/1/reviews"
```

### レビュー投稿（認証必要）
```bash
curl -X POST "http://localhost:8080/businesses/1/reviews" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"rating": 5, "text": "とても美味しかったです！"}'
```

### ユーザー閲覧履歴取得（認証必要）
```bash
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/logs/user/1/history"
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
│   ├── auth/                    # 認証サービス
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Dockerfile
│   │   ├── database/
│   │   ├── models/
│   │   └── handlers/
│   ├── business/                # ビジネスサービス
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Dockerfile
│   │   ├── database/
│   │   ├── models/
│   │   └── handlers/
│   ├── review/                  # レビューサービス
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Dockerfile
│   │   ├── database/
│   │   ├── models/
│   │   └── handlers/
│   └── logging/                 # ログサービス
│       ├── main.go
│       ├── go.mod
│       ├── go.sum
│       ├── Dockerfile
│       ├── cassandra/
│       ├── models/
│       └── handlers/
```

## セキュリティ機能

### ネットワーク分離
- **Frontend Network**: APIゲートウェイのみ外部アクセス可能
- **Backend Network**: 内部サービス間通信専用（external: false）

### 認証・認可
- **JWT認証**: ステートレス認証トークン
- **パスワードハッシュ化**: bcrypt使用
- **認証責務分離**: APIゲートウェイで一元管理

### API保護
- **認証必須エンドポイント**: レビュー投稿、履歴取得等
- **オプション認証**: パブリックエンドポイントでのログ記録

## サンプルデータ

`sample_data.sql`には以下のサンプルデータが含まれています：

- 5名のユーザー（パスワード付き）
- 8件のビジネス（ラーメン店、カフェ、寿司店、イタリアン、居酒屋、焼肉店、ベーカリー、中華料理）
- 21件のレビュー

