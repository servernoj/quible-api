package v1_test

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/stretchr/testify/assert"
)

type TCExtraTest func(TCRequest, *httptest.ResponseRecorder) bool
type TCRequest struct {
	Method  string
	Path    string
	Args    []any
	Body    map[string]any
	Headers []any
	Params  map[string]any
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

func (scenario *TCData) GetRunner(testAPI humatest.TestAPI) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		assert := assert.New(t)
		var state any
		// pre-hook (per-subtest initialization)
		if scenario.PreHook != nil {
			state = scenario.PreHook(t)
		}
		response := testAPI.Do(
			scenario.Request.Method,
			scenario.Request.Path,
			scenario.Request.Args...,
		)
		// response status
		assert.EqualValues(scenario.Response.Status, response.Code, "response status should match the expectation")
		// error code in case of error
		if scenario.Response.ErrorCode != nil {
			assert.Contains(
				response.Body.String(),
				strconv.Itoa(int(*scenario.Response.ErrorCode)),
				"error code should match expectation",
			)
		}
		// extra tests (if present)
		for _, fn := range scenario.ExtraTests {
			assert.True(
				fn(scenario.Request, response),
			)
		}
		// post-hook (post execution assertion)
		if scenario.PostHook != nil {
			scenario.PostHook(t, state)
		}
	}
}
