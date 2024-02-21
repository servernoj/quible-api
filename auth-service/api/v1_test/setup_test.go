package v1_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	POSTGRES_PASSWORD = "secret"
	POSTGRES_USER     = "user"
	POSTGRES_DB       = "dbname"
)

type TestSuite struct {
	suite.Suite
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (suite *TestSuite) SetupTest() {
	var err error
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
	boil.SetDB(db)
}
func (suite *TestSuite) TearDownTest() {
	if db, ok := boil.GetDB().(*sql.DB); ok {
		db.Close()
	}
	if err := suite.pool.Purge(suite.resource); err != nil {
		suite.T().Fatalf("Could not purge resource: %s", err)
	}
}

func (suite *TestSuite) TestOne() {
	assert := assert.New(suite.T())
	assert.Equal(true, true, "true is true")
}

func (suite *TestSuite) TestTwo() {
	assert := assert.New(suite.T())
	assert.Equal(false, false, "false is false")
}

func TestV1(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	suite.Run(t, &TestSuite{})
}
