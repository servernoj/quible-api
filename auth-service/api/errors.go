package api

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

type ErrorReporter interface {
	NewError(int, string, ...error) huma.StatusError
}

const ErrStatusGain = 10000

type ErrorCodeConstraints interface {
	~int
	String() string
}

// Custom error response
// implements https://pkg.go.dev/github.com/danielgtaylor/huma/v2#StatusError
type ErrorResponse struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (er *ErrorResponse) GetStatus() int {
	if er.Status > 0 {
		return er.Status
	}
	return er.Code / ErrStatusGain
}
func (er *ErrorResponse) Error() string {
	return er.Message
}

type ErrorMap[EC ErrorCodeConstraints] map[EC]string

func (errorMap ErrorMap[EC]) GetErrorResponse(errorCode EC, errs ...error) (er *ErrorResponse) {
	defer func() {
		log.Error().Msgf("Error[%d] %s", er.Code, er.Message)
		if len(errs) > 0 {
			log.Error().Errs("errors", errs).Send()
		}
	}()
	if msg, ok := errorMap[errorCode]; ok {
		er = &ErrorResponse{
			Code:    int(errorCode),
			Message: msg,
		}
		return
	}
	er = &ErrorResponse{
		Status:  http.StatusInternalServerError,
		Code:    int(errorCode),
		Message: fmt.Sprintf("Unable to process unmapped error code %d", errorCode),
	}
	return
}
