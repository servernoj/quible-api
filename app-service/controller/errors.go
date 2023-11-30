package controller

import (
	"net/http"

	c "github.com/quible-io/quible-api/lib/controller"
)

const ErrStatusGain = 10000
const ErrServiceId = 2000

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

// TODO: Complete the mapping
var errorMap = c.ErrorMap[ErrorCode]{
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
	http.StatusFailedDependency: {
		Err424_UnknownError: "unknown error",
	},
	// 500
	http.StatusInternalServerError: {
		Err500_UnknownError: "internal server error",
	},
}

var (
	SendError     = errorMap.SendError
	GetErrorCodes = errorMap.GetErrorCodes
)
