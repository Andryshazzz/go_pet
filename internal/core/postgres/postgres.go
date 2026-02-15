package postgres

import (
	"context"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, host, port, user, password, dbname string) (*pgxpool.Pool, error) {

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   dbname,
	}

	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	cfg, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctxTimeout, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctxTimeout); err != nil {
		return nil, err
	}

	return pool, nil
}
