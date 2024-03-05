package v1_test

import (
	_ "embed"
	"net/http"

	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/store"
)

func (suite *TestCases) TestPasswordReset() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Load environment variables
	env.Setup()
	// 3. Define test scenarios
	testCases := TCScenarios{}
	// 4. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPost, "/api/v1/user/password-reset"))
	}
}
