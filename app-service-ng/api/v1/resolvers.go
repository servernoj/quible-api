package v1

import (
	"github.com/danielgtaylor/huma/v2"
)

type _resolved struct {
	resolved bool
}

// -- Authorization header containing Bearer access token. Injects `UserId` into `input` struct
type AuthorizationHeaderResolver struct {
	_resolved
	Authorization string `header:"authorization"`
	UserId        string
}

func (f *AuthorizationHeaderResolver) Resolve(ctx huma.Context) (errs []error) {
	return
}
