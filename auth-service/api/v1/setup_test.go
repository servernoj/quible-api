package v1_test

import (
	_ "embed"
	"net/http/httptest"
	"os"
	"testing"

	srvAPI "github.com/quible-io/quible-api/auth-service/api"
	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

//go:embed TestData/users.csv
var users_as_csv string

type TestCases struct {
	libAPI.TestSuite
	libAPI.ServiceAPI
}
type TCExtraTest func(TCRequest, *httptest.ResponseRecorder) bool
type TCRequest struct {
	Body map[string]any
}
type TCResponse struct {
	Status    int
	ErrorCode *v1.ErrorCode
}
type TCData struct {
	Description string
	Request     TCRequest
	Response    TCResponse
	ExtraTests  []TCExtraTest
	PreHook     func(*testing.T) any
	PostHook    func(*testing.T, any)
}
type TCScenarios map[string]TCData

// This is the only test function being called by `go test ./...` It takes advantage of `testify/suite` package
// to initialize a test suite containing (implementing) `SetupTest` and `TearDownTest` methods that are automatically
// called before and after "each test". The "each test" term defines methods in the `TestCases` that have names started with `Test`,
// for example `TestUserLogin`.
func TestRunner(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	implementation := v1.New()
	suite.Run(
		t,
		&TestCases{
			ServiceAPI: implementation,
			TestSuite: libAPI.NewTestSuite(
				t,
				implementation,
				srvAPI.Title,
				libAPI.VersionConfig{
					Tag:    "v1",
					SemVer: "1.0.0",
				},
			),
		},
	)
}
