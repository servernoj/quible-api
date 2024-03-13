package v1_test

import (
	_ "embed"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/gin-gonic/gin"
	srvAPI "github.com/quible-io/quible-api/app-service/api"
	v1 "github.com/quible-io/quible-api/app-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed TestData/users.csv
var UsersCSV string

//go:embed TestData/chats.csv
var ChatsCSV string

//go:embed TestData/chat-user.csv
var ChatUserCSV string

type tlogWriter struct {
	t *testing.T
}

func (lw *tlogWriter) Write(p []byte) (n int, err error) {
	lw.t.Helper()
	lw.t.Logf((string)(p))
	return len(p), nil
}

type (
	TestCases struct {
		suite.TestSuite
		humatest.TestAPI
		libAPI.ServiceAPI
	}
)

func TestRunner(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: &tlogWriter{t}})
	gin.SetMode(gin.ReleaseMode)
	serviceAPI := v1.NewServiceAPI(
		v1.WithDeps(
			libAPI.NewDeps(
				make(map[string]any),
			),
		),
	)
	api := srvAPI.Setup(
		serviceAPI,
		gin.Default(),
		libAPI.VersionConfig{},
	)

	testAPI := suite.NewTestAPI(t, api)
	suite.RunSuite(
		t,
		&TestCases{
			ServiceAPI: serviceAPI,
			TestSuite: suite.TestSuite{
				DBStore: suite.NewDBs(),
			},
			TestAPI: testAPI,
		},
		true,
	)
}
