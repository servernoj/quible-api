package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConfig is configuration wrapper for our database setting
type DatabaseConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	DBName   string
}

// DSN will get datasource name of the database configuration
func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		db.Username, db.Password, db.Hostname, db.Port, db.DBName)
}

// NewDBPool will create pool connection to database
func NewDBPool(dbConfig DatabaseConfig) (*pgxpool.Pool, func(), error) {

	f := func() {}

	// create pgx connection pool
	pool, err := pgxpool.New(context.Background(), dbConfig.DSN())

	// return nil to connection and return error if error occur
	if err != nil {
		return nil, f, errors.New("database connection error")
	}

	// validateDBPool
	err = validateDBPool(pool)

	// return nil and error if error occur
	if err != nil {
		return nil, f, err
	}

	// return connection pool and inline function to close/ clear the pool if not used.
	// return nil for the error since there should be no error to this point
	return pool, func() { pool.Close() }, nil
}

// validateDBPool will pings the database and logs the current user and database
func validateDBPool(pool *pgxpool.Pool) error {
	// tried to ping connection
	err := pool.Ping(context.Background())

	// return error if error found
	if err != nil {
		return errors.New("database connection error")
	}

	var (
		currentDatabase string
		currentUser     string
		dbVersion       string
	)

	// Lets try to get db system info
	sqlStatement := `select current_database(), current_user, version();`
	row := pool.QueryRow(context.Background(), sqlStatement)
	err = row.Scan(&currentDatabase, &currentUser, &dbVersion)

	switch {
	case err == sql.ErrNoRows:
		return errors.New("no rows were returned")
	case err != nil:
		return errors.New("database connection error")
	default:
		log.Printf("database version: %s\n", dbVersion)
		log.Printf("database user: %s\n", currentUser)
		log.Printf("database name: %s\n", currentDatabase)
	}

	return nil
}
