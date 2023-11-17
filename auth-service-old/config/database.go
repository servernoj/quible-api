package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBPool will create pool connection to database
func NewDBPool(dsn string) (*pgxpool.Pool, error) {

	ctx := context.Background()

	// create pgx connection pool
	pool, err := pgxpool.New(ctx, dsn)

	// return nil to connection and return error if error occur
	if err != nil {
		return nil, fmt.Errorf("opening DB driver: %w", err)
	}

	// return nil and error if error occur
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging DB: %w", err)
	}

	return pool, nil
}
