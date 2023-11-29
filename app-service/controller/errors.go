package controller

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

const ErrStatusGain = 10000
const ErrServiceId = 2000

// -- specific case for `app-service`

//go:generate stringer -type=ErrorCode
type ErrorCode int

const (
	Err400_Shift = ErrStatusGain*http.StatusBadRequest + ErrServiceId
	Err401_Shift = ErrStatusGain*http.StatusUnauthorized + ErrServiceId
	Err404_Shift = ErrStatusGain*http.StatusNotFound + ErrServiceId
	Err424_Shift = ErrStatusGain*http.StatusFailedDependency + ErrServiceId
	Err500_Shift = ErrStatusGain*http.StatusInternalServerError + ErrServiceId
)

const (
	Err400_UnknownError ErrorCode = Err400_Shift + iota + 1
	Err400_MalformedJSON
	Err400_InvalidRequestBody
)
const (
	Err401_UnknownError ErrorCode = Err401_Shift + iota + 1
	Err401_UserIdNotFound
	Err401_UserNotFound
)
const (
	Err404_UnknownError ErrorCode = Err404_Shift + iota + 1
)
const (
	Err424_UnknownError ErrorCode = Err424_Shift + iota + 1
)
const (
	Err500_UnknownError ErrorCode = Err500_Shift + iota + 1
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorMap map[int]map[ErrorCode]string

// TODO: Complete the mapping
var errorMap = ErrorMap{
	// 400
	http.StatusBadRequest: {
		Err400_UnknownError:       "unknown error",
		Err400_MalformedJSON:      "malformed JSON request",
		Err400_InvalidRequestBody: "invalid request body",
	},
	// 401
	http.StatusUnauthorized: {
		Err401_UnknownError:   "unknown error",
		Err401_UserIdNotFound: "userId not present",
		Err401_UserNotFound:   "user not found",
	},
	// 404
	http.StatusNotFound: {
		Err404_UnknownError: "unknown error",
	},
	// 424
	http.StatusTooManyRequests: {
		Err424_UnknownError: "unknown error",
	},
	// 500
	http.StatusInternalServerError: {
		Err500_UnknownError: "internal server error",
	},
}

func SendError(c *gin.Context, status int, code ErrorCode) {
	if _, ok := errorMap[status][code]; !ok {
		status = http.StatusInternalServerError
		code = Err500_UnknownError
	}
	c.JSON(
		status,
		ErrorResponse{
			Code:    int(code),
			Message: errorMap[status][code],
		},
	)
	c.Abort()
}

func GetErrorCodes(c *gin.Context) {
	statuses := sort.IntSlice(make([]int, len(errorMap)))
	idx := 0
	for status := range errorMap {
		statuses[idx] = status
		idx++
	}
	sort.Sort(statuses)
	builder := strings.Builder{}
	builder.WriteString("<h2>Error codes</h2>")
	for _, httpStatus := range statuses {
		statusMap := errorMap[httpStatus]
		builder.WriteString(fmt.Sprintf("<h3>Status code: %d</h3>", httpStatus))
		builder.WriteString("<pre>")

		// sort statusMap keys
		errorCodes := sort.IntSlice(make([]int, len(statusMap)))
		idx := 0
		for errorCode := range statusMap {
			errorCodes[idx] = int(errorCode)
			idx++
		}
		sort.Sort(errorCodes)
		for _, errorCode := range errorCodes {
			errorCode := ErrorCode(errorCode)
			message := statusMap[errorCode]
			builder.WriteString(
				fmt.Sprintf("%-50s%-20d%s\n", errorCode.String(), errorCode, message),
			)
		}
		builder.WriteString("</pre>")
	}
	c.Writer.Header().Add("content-type", gin.MIMEHTML)
	c.String(http.StatusOK, builder.String())
}
