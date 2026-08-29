package postgrespool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool defines the interface for database operations.
// It abstracts pgxpool.Pool to allow mocking in tests.
//
// All methods accept context for cancellation and timeouts.
type Pool interface {
	// Query executes a SELECT query and returns rows.
	// The caller must close the returned rows.
	//
	// Usage:
	//   rows, err := pool.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
	//   if err != nil { ... }
	//   defer rows.Close()
	//   for rows.Next() { ... }
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	// QueryRow executes a SELECT query expected to return at most one row.
	//
	// Usage:
	//   var name string
	//   err := pool.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", userID).Scan(&name)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE).
	//
	// Usage:
	//   tag, err := pool.Exec(ctx, "UPDATE users SET active = $1 WHERE id = $2", true, userID)
	//   if tag.RowsAffected() == 0 { ... }
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

	// Close releases all connections in the pool.
	// Should be called once on application shutdown.
	Close()

	// OpTimeout returns the configured operation timeout.
	// Repositories should use it to create contexts:
	//   ctx, cancel := context.WithTimeout(ctx, pool.OpTimeout())
	OpTimeout() time.Duration
}

// ConnectionPool wraps pgxpool.Pool with an operation timeout.
// It implements the Pool interface.
//
// Usage:
//
//	pool, err := postgrespool.NewConnectionPool(ctx, postgrespool.NewConfigMust())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pool.Close()
//
//	rows, err := pool.Query(ctx, "SELECT * FROM users")
type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

// NewConnectionPool creates a new PostgreSQL connection pool.
// It establishes the connection, pings the database to verify connectivity,
// and returns a ready-to-use pool.
//
// The connection string is built from the config:
//
//	postgres://user:password@host:port/database?sslmode=disable
//
// Usage:
//
//	cfg := postgrespool.NewConfigMust()
//	pool, err := postgrespool.NewConnectionPool(ctx, cfg)
//	if err != nil {
//	    log.Fatal("failed to connect to database:", err)
//	}
//	defer pool.Close()
func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		return nil, fmt.Errorf("Parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)

	if err != nil {
		return nil, fmt.Errorf("Сreate pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("Pgxpool ping: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

// OpTimeout returns the configured operation timeout duration.
//
// Usage in repository:
//
//	ctx, cancel := context.WithTimeout(ctx, pool.OpTimeout())
//	defer cancel()
//	rows, err := pool.Query(ctx, "SELECT ...")
func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
