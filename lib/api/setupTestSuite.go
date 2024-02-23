package api

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/quible-io/quible-api/lib/migrations"
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
	TestAPI  humatest.TestAPI
}

type MyTB struct {
	disableLogging bool
	*testing.T
}

func (tb *MyTB) Log(args ...any) {
	if !tb.disableLogging {
		tb.T.Log(args...)
	}
}
func (tb *MyTB) Logf(format string, args ...any) {
	if !tb.disableLogging {
		tb.T.Logf(format, args...)
	}
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
	// 4. Perform DB migrations
	goose.SetLogger(goose.NopLogger())
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

func NewTestSuite[Impl ErrorReporter](t *testing.T, vc VersionConfig) TestSuite {
	// 1. Mimic error response from the actual API
	var implValue Impl
	overrideHumaNewError(implValue)
	// 2. Create new test API
	myTB := &MyTB{
		T:              t,
		disableLogging: true,
	}
	_, api := humatest.New(myTB)
	// 3. Register operations from the implemented API to the test API
	implType := reflect.TypeOf(&implValue)
	args := []reflect.Value{
		reflect.ValueOf(&implValue),
		reflect.ValueOf(api),
		reflect.ValueOf(vc),
	}
	for i := 0; i < implType.NumMethod(); i++ {
		m := implType.Method(i)
		if strings.HasPrefix(m.Name, "Register") && len(m.Name) > 8 {
			m.Func.Call(args)
		}
	}
	return TestSuite{
		TestAPI: api,
	}
}
