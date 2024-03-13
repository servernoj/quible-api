package suite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/quible-io/quible-api/lib/migrations"
	"github.com/rs/zerolog/log"
)

const (
	POSTGRES_PASSWORD = "secret"
	POSTGRES_USER     = "user"
	POSTGRES_DB       = "dbname"
)

type Suite interface {
	SetupTest(*testing.T) *TestData
	TearDownTest(*testing.T, *TestData)
}
type DBStore interface {
	StoreDB(string, *sql.DB)
	RetrieveDB(string) *sql.DB
}
type SuiteStore interface {
	Suite
	DBStore
}

type TestSuite struct {
	Suite
	DBStore
}

type TestData struct {
	*sql.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (ts *TestSuite) SetupTest(t *testing.T) *TestData {
	testData := TestData{}
	var err error
	// 1. Initialize pool
	testData.pool, err = dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not construct pool: %s", err)
	}
	err = testData.pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}
	// 2. Initialize resource
	testData.resource, err = testData.pool.RunWithOptions(&dockertest.RunOptions{
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
		t.Fatalf("Could not start resource: %s", err)
	}
	// 3. Connect to DB and initialize SQLBoiler
	hostAndPort := testData.resource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		POSTGRES_USER,
		POSTGRES_PASSWORD,
		hostAndPort,
		POSTGRES_DB,
	)
	log.Info().Msgf("Connecting to database...")
	if err := testData.resource.Expire(120); err != nil {
		t.Fatalf("Could not set hard kill timer for the container: %s", err)
	}
	testData.pool.MaxWait = 120 * time.Second
	var db *sql.DB
	if err = testData.pool.Retry(
		func() error {
			db, err = sql.Open("pgx", dsn)
			if err != nil {
				return err
			}
			return db.Ping()
		},
	); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	log.Info().Msg("Database connected")
	testData.DB = db
	// 4. Perform DB migrations
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations.FS)
	ctx := context.Background()
	if err := goose.RunContext(ctx, "up", db, "."); err != nil {
		t.Fatalf("Could not migrate DB: %s", err)
	}
	version, _ := goose.GetDBVersionContext(ctx, db)
	log.Info().Msgf("DB migrated up to %d", version)
	// 5. Return test specific data
	return &testData
}

func (ts *TestSuite) TearDownTest(t *testing.T, testData *TestData) {
	if err := testData.DB.Close(); err != nil {
		t.Fatalf("Could not close DB handle: %s", err)
	}
	if err := testData.pool.Purge(testData.resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}
