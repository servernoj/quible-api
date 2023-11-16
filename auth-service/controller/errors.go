package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorMap map[int]map[int]string

// TODO: Complete the mapping
var errorMap = ErrorMap{
	// 207
	http.StatusMultiStatus: {
		0001: "Partial Deletion: Some data remains undeleted",
	},
	// 400
	http.StatusBadRequest: {
		0001: "email is not registered",
		0002: "wrong password",
		0003: "invalid email address format",
		0004: "invalid username format",
		0005: "username already exists",
		0006: "password doesn't match regex",
		0007: "Bad Request: Invalid phone number format",
	},
	// 401
	http.StatusUnauthorized: {
		0001: "Unauthorized: Invalid credentials provided",
		0002: "Unauthorized: Invalid credentials provided",
	},
	// 403
	http.StatusForbidden: {
		0001: "Forbidden: Insufficient permissions for deletion",
		0002: "Forbidden: Insufficient permissions for phone number edit",
	},
	// 404
	http.StatusNotFound: {
		0001: "Not Found: Player stats not Available",
		0002: "Not Found: Player stats not Available",
		0003: "Not Found: User or phone number not found",
		0004: "Not Found: Account already deleted or does not exist",
	},
	// 429
	http.StatusTooManyRequests: {
		0001: "Limit Exceeded: Edit requests reached limit",
	},
	// 500
	http.StatusInternalServerError: {
		0001: "Internal Server Error: Unexpected issue during deletion",
		0002: "Internal Server Error: Unexpected issue during phone number edit",
		0003: "Internal Server Error: Unexpected issue occurred",
		0004: "Internal Server Error: Unexpected issue occurred",
		9999: "Unknown error",
	},
	// 503
	http.StatusServiceUnavailable: {
		0001: "Service Unavailable: Database issue during deletion",
		0002: "Service Unavailable: Database issue during phone number edit",
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500_9999,
			Message: errorMap[500][9999],
		})
	}
}
