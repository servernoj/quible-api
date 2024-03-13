package v1

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email/postmark"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

//go:embed serviceDescription.md
var ServiceDescription string

type WithOption func(*VersionedImpl)

func WithDeps(deps libAPI.Deps) WithOption {
	return func(vi *VersionedImpl) {
		vi.Deps = deps
	}
}

func NewServiceAPI(opts ...WithOption) libAPI.ServiceAPI {
	impl := &VersionedImpl{
		Deps: libAPI.NewDeps(
			map[string]any{
				"db":     boil.GetDB(),
				"mailer": postmark.NewClient(),
			},
		),
	}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}

type VersionedImpl struct {
	libAPI.Deps
}

func (impl VersionedImpl) NewError(status int, message string, errs ...error) huma.StatusError {
	if status == http.StatusUnprocessableEntity && message == "validation failed" {
		locationToErrorCode := map[string]ErrorCode{
			"header.authorization": Err401_InvalidAccessToken,
			"auth-service":         Err401_AuthServiceError,
			"db.users":             Err401_UserNotFound,
		}
		for i := 0; i < len(errs); i++ {
			if converted, ok := errs[i].(huma.ErrorDetailer); ok {
				location := converted.ErrorDetail().Location
				for key, errorCode := range locationToErrorCode {
					if strings.Contains(location, key) {
						return ErrorMap.GetErrorResponse(errorCode)
					}
				}
			}
		}
		return ErrorMap.GetErrorResponse(Err400_InvalidRequest)
	}
	return ErrorMap.GetErrorResponse(Err500_UnknownHumaError, errs...)
}
