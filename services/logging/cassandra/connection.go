package cassandra

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

var Session gocqlx.Session

func Connect() error {
	hosts := os.Getenv("CASSANDRA_HOSTS")
	if hosts == "" {
		hosts = "localhost:9042"
	}

	cluster := gocql.NewCluster(hosts)
	//- Quorum: 読み書き時に過半数のレプリカから応答を要求
	//- データの一貫性と可用性のバランスが良い設定
	//- 例：3レプリカなら2つのノードから応答が必要
	cluster.Consistency = gocql.Quorum
	//- CassandraのCQLプロトコルバージョン4を使用
	//- Cassandra 2.2以降でサポート
	//- より効率的な通信とパフォーマンス向上
	cluster.ProtoVersion = 4
	//- 初回接続時のタイムアウト: 10秒
	//- ノードへの接続確立の待ち時間
	cluster.ConnectTimeout = time.Second * 10
	//- クエリ実行のタイムアウト: 10秒
	//- 個々のクエリの最大実行時間
	cluster.Timeout = time.Second * 10
	//- 失敗時の再試行回数: 3回
	//- ネットワークエラーや一時的な障害に対応
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return fmt.Errorf("failed to connect to Cassandra: %w", err)
	}

	Session = session
	log.Println("Connected to Cassandra successfully")
	return nil
}

func Close() {
	if Session.Session != nil {
		Session.Close()
		log.Println("Cassandra connection closed")
	}
}

func AutoMigrate() error {
	log.Println("Running Cassandra migrations...")

	// Create keyspace
	//replication_factor: 1
	//	- SimpleStrategy: 単一データセンター向け
	//	- レプリカ数1: 開発環境向け（本番では3以上推奨）
	createKeyspace := `
		CREATE KEYSPACE IF NOT EXISTS yelp_logs WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}
	`

	if err := Session.ExecStmt(createKeyspace); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	// Create review_view_logs table (with keyspace prefix)
	createTable := `
		CREATE TABLE IF NOT EXISTS yelp_logs.review_view_logs (
			user_id INT,
			business_id INT,
			review_id INT,
			viewed_at TIMESTAMP,
			ip_address TEXT,
			user_agent TEXT,
			PRIMARY KEY (user_id, viewed_at, business_id)
		) WITH CLUSTERING ORDER BY (viewed_at DESC)
	`

	if err := Session.ExecStmt(createTable); err != nil {
		return fmt.Errorf("failed to create review_view_logs table: %w", err)
	}

	// Create schema_migrations table for tracking (with keyspace prefix)
	createMigrationTable := `
		CREATE TABLE IF NOT EXISTS yelp_logs.schema_migrations (
			version TEXT PRIMARY KEY,
			executed_at TIMESTAMP,
			checksum TEXT
		)
	`

	if err := Session.ExecStmt(createMigrationTable); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	log.Println("Cassandra migrations completed successfully")
	return nil
}
