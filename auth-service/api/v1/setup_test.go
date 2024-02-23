package v1_test

import (
	"os"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type TestCases struct {
	libAPI.TestSuite
}

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
}
type TCScenarios map[string]TCData

// This is the only test function being called by `go test ./...` It takes advantage of `testify/suite` package
// to initialize a test suite containing (implementing) `SetupTest` and `TearDownTest` methods that are automatically
// called before and after "each test". The "each test" term defines methods in the suit that have names started with `Test`,
// for example `TestUserLogin`.
func TestRunner(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	suite.Run(
		t,
		&TestCases{
			TestSuite: libAPI.NewTestSuite[v1.VersionedImpl](
				t,
				libAPI.VersionConfig{
					Tag:    "v1",
					SemVer: "1.0.0",
				},
			),
		},
	)
}
