package v1

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type VersionedImpl struct{}

type UserSimplified struct {
	ID       string `json:"id" doc:"user ID (UUID)"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

type UserProfile struct {
	ID       string  `json:"id" doc:"user ID (UUID)"`
	FullName string  `json:"full_name"`
	Image    *string `json:"image" doc:"Profile image data URL"`
}

type ImageData struct {
	ContentType   string `json:"contentType"`
	BinaryContent []byte `json:"data"`
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
