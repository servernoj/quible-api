package v1

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/email/postmark"
)

func New() libAPI.ServiceAPI {
	return &VersionedImpl{
		EmailSender: postmark.NewClient(),
	}
}

type VersionedImpl struct {
	email.EmailSender
}

func (impl *VersionedImpl) Register(api huma.API, vc libAPI.VersionConfig) {
	implType := reflect.TypeOf(impl)
	args := []reflect.Value{
		reflect.ValueOf(impl),
		reflect.ValueOf(api),
		reflect.ValueOf(vc),
	}
	for i := 0; i < implType.NumMethod(); i++ {
		m := implType.Method(i)
		if strings.HasPrefix(m.Name, "Register") && len(m.Name) > 8 {
			m.Func.Call(args)
		}
	}
}

func (impl VersionedImpl) NewError(status int, message string, errs ...error) huma.StatusError {
	if status == http.StatusUnprocessableEntity && message == "validation failed" {
		locationToErrorCode := map[string]ErrorCode{
			"body.email":           Err400_InvalidEmailFormat,
			"body.phone":           Err400_InvalidPhoneFormat,
			"body.password":        Err400_UnsatisfactoryPassword,
			"body.confirmPassword": Err400_UnsatisfactoryConfirmPassword,
			"body.token":           Err400_InvalidOrMalformedToken,
			"header.authorization": Err401_InvalidAccessToken,
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
