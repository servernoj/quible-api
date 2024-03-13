package api

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type TCExtraTest func(TCRequest, *httptest.ResponseRecorder) bool
type TCRequest struct {
	Args   []any
	Params map[string]any
}
type TCResponse struct {
	Status    int
	ErrorCode *int
}
type TCEnv map[string]string
type TCData struct {
	Description string
	Envs        TCEnv
	Request     TCRequest
	Response    TCResponse
	ExtraTests  []TCExtraTest
	PreHook     func(*testing.T) any
	PostHook    func(*testing.T, any)
}
type TCScenarios map[string]TCScenario

type TCScenario func(*testing.T) TCData

func (scenario TCScenario) GetRunner(testAPI humatest.TestAPI, method string, pathFormat string, pathParamsKeys ...string) func(*testing.T) {
	return func(t *testing.T) {
		tcData := scenario(t)
		assert := assert.New(t)
		// set/unset environment variables
		for k, v := range tcData.Envs {
			os.Setenv(k, v)
		}
		t.Cleanup(
			func() {
				for k := range tcData.Envs {
					os.Unsetenv(k)
				}
			},
		)
		var state any
		// pre-hook (per-subtest initialization)
		if tcData.PreHook != nil {
			state = tcData.PreHook(t)
		}
		// inflate pathParams into pathArgs via tcData.Params map
		pathArgs := make([]any, len(pathParamsKeys))
		for idx, key := range pathParamsKeys {
			pathArgs[idx] = tcData.Request.Params[key]
		}
		response := testAPI.Do(
			method,
			fmt.Sprintf("/api"+pathFormat, pathArgs...),
			tcData.Request.Args...,
		)
		// response status
		assert.EqualValues(tcData.Response.Status, response.Code, "response status should match the expectation")
		// error code in case of error
		if tcData.Response.ErrorCode != nil {
			assert.Contains(
				response.Body.String(),
				strconv.Itoa(*tcData.Response.ErrorCode),
				"error code should match expectation",
			)
		}
		// extra tests (if present)
		for _, fn := range tcData.ExtraTests {
			assert.True(
				fn(tcData.Request, response),
			)
		}
		// post-hook (post execution assertion)
		if tcData.PostHook != nil {
			tcData.PostHook(t, state)
		}
		log.Info().Msg("Done")
	}
}
