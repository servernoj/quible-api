package v1

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

func TestV1(t *testing.T) {
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
