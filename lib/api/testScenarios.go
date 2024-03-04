package api

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

type TCExtraTest func(TCRequest, *httptest.ResponseRecorder) bool
type TCRequest struct {
	Args   []any
	Params map[string]any
}
type TCResponse[EC ErrorCodeConstraints] struct {
	Status    int
	ErrorCode *EC
}
type TCData[EC ErrorCodeConstraints] struct {
	Description string
	Request     TCRequest
	Response    TCResponse[EC]
	ExtraTests  []TCExtraTest
	PreHook     func(*testing.T) any
	PostHook    func(*testing.T, any)
}
type TCScenarios[EC ErrorCodeConstraints] map[string]TCData[EC]

func (scenario *TCData[EC]) GetRunner(testAPI humatest.TestAPI, method string, pathFormat string, pathArgs ...any) func(*testing.T) {
	return func(t *testing.T) {
		assert := assert.New(t)
		var state any
		// pre-hook (per-subtest initialization)
		if scenario.PreHook != nil {
			state = scenario.PreHook(t)
		}
		response := testAPI.Do(
			method,
			fmt.Sprintf(pathFormat, pathArgs...),
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
