package v1_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/quible-io/quible-api/cmd/migrations"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	POSTGRES_PASSWORD = "secret"
	POSTGRES_USER     = "user"
	POSTGRES_DB       = "dbname"
)

type DebugWriter struct{}

func (DebugWriter) Write(p []byte) (n int, err error) {
	log.Debug().Msg(strings.TrimSpace(string(p)))
	return len(p), nil
}

type TestSuite struct {
	suite.Suite
	pool     *dockertest.Pool
	resource *dockertest.Resource
	ctx      context.Context
}

func (suite *TestSuite) SetupTest() {
	var err error
	// 0. Initialize context
	suite.ctx = context.Background()
	// 1. Initialize pool
	suite.pool, err = dockertest.NewPool("")
	if err != nil {
		suite.T().Fatalf("Could not construct pool: %s", err)
	}
	err = suite.pool.Client.Ping()
	if err != nil {
		suite.T().Fatalf("Could not connect to Docker: %s", err)
	}
	// 2. Initialize resource
	suite.resource, err = suite.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", POSTGRES_PASSWORD),
			fmt.Sprintf("POSTGRES_USER=%s", POSTGRES_USER),
			fmt.Sprintf("POSTGRES_DB=%s", POSTGRES_DB),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		suite.T().Fatalf("Could not start resource: %s", err)
	}
	// 3. Connect to DB and initialize SQLBoiler
	hostAndPort := suite.resource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		POSTGRES_USER,
		POSTGRES_PASSWORD,
		hostAndPort,
		POSTGRES_DB,
	)
	log.Info().Msgf("Connecting to database...")
	if err := suite.resource.Expire(120); err != nil {
		suite.T().Fatalf("Could not set hard kill timer for the container: %s", err)
	}
	suite.pool.MaxWait = 120 * time.Second
	var db *sql.DB
	if err = suite.pool.Retry(
		func() error {
			db, err = sql.Open("pgx", dsn)
			if err != nil {
				return err
			}
			return db.Ping()
		},
	); err != nil {
		suite.T().Fatalf("Could not connect to docker: %s", err)
	}
	log.Info().Msg("Database connected")
	// 4. Perform DB migrations
	goose.SetBaseFS(migrations.FS)
	if err := goose.Run("up", db, "."); err != nil {
		suite.T().Fatalf("Could not migrate DB: %s", err)
	}
	log.Info().Msg("DB migrated")
	// 5. Setup SQLBoiler
	boil.SetDB(db)
	boil.DebugMode = true
	boil.DebugWriter = new(DebugWriter)
}

func (suite *TestSuite) TearDownTest() {
	if db, ok := boil.GetDB().(*sql.DB); ok {
		db.Close()
	}
	if err := suite.pool.Purge(suite.resource); err != nil {
		suite.T().Fatalf("Could not purge resource: %s", err)
	}
}

func TestV1(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	suite.Run(t, &TestSuite{})
}
