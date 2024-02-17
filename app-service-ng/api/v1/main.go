package v1

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type VersionedImpl struct{}

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
		return ErrorMap.GetErrorResponse(Err400_InvalidRequestBody)
	}
	return ErrorMap.GetErrorResponse(Err500_UnknownHumaError, errs...)
}
