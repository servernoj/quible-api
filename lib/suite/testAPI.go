package suite

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

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

func NewTestAPI(t *testing.T, api huma.API) humatest.TestAPI {
	return humatest.Wrap(
		&MyTB{
			T:              t,
			disableLogging: true,
		},
		api,
	)
}
