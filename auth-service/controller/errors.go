package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const generalErrorCode = 9999

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorMap map[int]map[int]string

// TODO: Complete the mapping
var errorMap = ErrorMap{
	// 207
	http.StatusMultiStatus: {
		1: "some data remains undeleted",
	},
	// 400
	http.StatusBadRequest: {
		1:   "email is not registered",
		2:   "wrong password",
		3:   "invalid email address format",
		4:   "invalid username format",
		5:   "user with such username exists",
		6:   "password doesn't match regex",
		7:   "invalid phone number format",
		8:   "user with such email exists",
		100: "invalid request body",
	},
	// 401
	http.StatusUnauthorized: {
		1: "invalid credentials provided",
		2: "authorization header missing",
		3: "authorization header is invalid",
		4: "no user found",
	},
	// 403
	http.StatusForbidden: {
		1: "insufficient permissions for deletion",
		2: "insufficient permissions for phone number edit",
	},
	// 404
	http.StatusNotFound: {
		1: "player stats not Available",
		2: "player stats not Available",
		3: "user or phone number not found",
		4: "account already deleted or does not exist",
	},
	// 429
	http.StatusTooManyRequests: {
		1: "edit requests reached limit",
	},
	// 500
	http.StatusInternalServerError: {
		1:                "unexpected issue during deletion",
		2:                "unexpected issue during phone number edit",
		3:                "unexpected issue occurred",
		4:                "unexpected issue during registration",
		generalErrorCode: "internal server error",
	},
	// 503
	http.StatusServiceUnavailable: {
		1: "Service Unavailable: Database issue during deletion",
		2: "Service Unavailable: Database issue during phone number edit",
	},
}

func SendError(c *gin.Context, httpStatus, messageId int) {
	if message, ok := errorMap[httpStatus][messageId]; ok {
		code, _ := strconv.Atoi(fmt.Sprintf("%03d%04d", httpStatus, messageId))
		errorResponse := ErrorResponse{
			Code:    code,
			Message: message,
		}
		c.JSON(httpStatus, errorResponse)
	} else {
		status := http.StatusInternalServerError
		code, _ := strconv.Atoi(
			fmt.Sprintf("%03d%04d", status, generalErrorCode),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    code,
			Message: errorMap[status][generalErrorCode],
		})
	}
}
