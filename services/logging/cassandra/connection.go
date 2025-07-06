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
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Timeout = time.Second * 10
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
	createKeyspace := `
		CREATE KEYSPACE IF NOT EXISTS yelp_logs WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}
	`

	if err := Session.ExecStmt(createKeyspace); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	// Note: gocql doesn't support USE statements, we specify keyspace in table names

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
