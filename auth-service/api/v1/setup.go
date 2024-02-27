package v1

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/email/postmark"
)

type WithOption func(*VersionedImpl)

func New() libAPI.ServiceAPI {
	return &VersionedImpl{
		EmailSender: postmark.NewClient(),
	}
}

type VersionedImpl struct {
	email.EmailSender
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

func (impl *VersionedImpl) SetEmailSender(emailSender email.EmailSender) {
	impl.EmailSender = emailSender
}
