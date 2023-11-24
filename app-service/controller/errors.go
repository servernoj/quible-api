package controller

/*
import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrTokenExpired              = errors.New("token expired")
	ErrTokenInvalidClaims        = errors.New("unable to process token claims")
	ErrTokenInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenInvalidType          = errors.New("invalid token type (access|refresh)")
	ErrTokenMissingUserId        = errors.New("unable to extract userId from token")
	ErrTokenMissingTokenId       = errors.New("unable to extract tokenId from token")
)

const ErrStatusGain = 10_000
const ErrServiceIdGain = 1_000

// -- specific case for `app-service`
const ErrServiceId = 1

//nogo:generate stringer -type=ErrorCode
type ErrorCode int

const (
	Err207_Shift = ErrStatusGain*http.StatusMultiStatus + ErrServiceIdGain*ErrServiceId
	Err400_Shift = ErrStatusGain*http.StatusBadRequest + ErrServiceIdGain*ErrServiceId
	Err401_Shift = ErrStatusGain*http.StatusUnauthorized + ErrServiceIdGain*ErrServiceId
	Err403_Shift = ErrStatusGain*http.StatusForbidden + ErrServiceIdGain*ErrServiceId
	Err404_Shift = ErrStatusGain*http.StatusNotFound + ErrServiceIdGain*ErrServiceId
	Err429_Shift = ErrStatusGain*http.StatusTooManyRequests + ErrServiceIdGain*ErrServiceId
	Err500_Shift = ErrStatusGain*http.StatusInternalServerError + ErrServiceIdGain*ErrServiceId
	Err503_Shift = ErrStatusGain*http.StatusServiceUnavailable + ErrServiceIdGain*ErrServiceId
)

const (
	Err207_SomeDataUndeleted ErrorCode = Err207_Shift + iota + 1
)
const (
	Err400_EmailNotRegistered ErrorCode = Err400_Shift + iota + 1
	Err400_InvalidEmailFormat
	Err400_InvalidUsernameFormat
	Err400_InvalidPhoneFormat
	Err400_UserWithUsernameExists
	Err400_UserWithEmailExists
	Err400_IsufficientPasswordComplexity
	Err400_MalformedJSON
	Err400_InvalidRequestBody
)

const (
	Err401_InvalidCredentials ErrorCode = Err401_Shift + iota + 1
	Err401_AuthorizationHeaderMissing
	Err401_AuthorizationHeaderInvalid
	Err401_AuthorizationExpired
	Err401_InvalidRefreshToken
	Err401_UserNotFound
)

const (
	Err403_CannotToDelete ErrorCode = Err403_Shift + iota + 1
	Err403_CannotEditPhone
)
const (
	Err404_PlayerStatsNotFound ErrorCode = Err404_Shift + iota + 1
	Err404_UserOrPhoneNotFound
	Err404_AccountNotFound
)
const (
	Err429_EditRequestTimedOut ErrorCode = Err429_Shift + iota + 1
)
const (
	Err500_UnknownError ErrorCode = Err500_Shift + iota + 1
	Err500_UnableToDelete
	Err500_UnableToEditPhone
	Err500_UnableToRegister
	Err500_UnableToGenerateToken
)
const (
	Err503_DataBaseOnDelete ErrorCode = Err503_Shift + iota + 1
	Err503_DataBaseOnPhoneEdit
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorMap map[int]map[ErrorCode]string

// TODO: Complete the mapping
var errorMap = ErrorMap{
	// 207
	http.StatusMultiStatus: {
		Err207_SomeDataUndeleted: "some data remains undeleted",
	},
	// 400
	http.StatusBadRequest: {
		Err400_EmailNotRegistered:            "email is not registered",
		Err400_InvalidEmailFormat:            "invalid email address format",
		Err400_InvalidUsernameFormat:         "invalid username format",
		Err400_InvalidPhoneFormat:            "invalid phone number format",
		Err400_UserWithUsernameExists:        "user with such username exists",
		Err400_IsufficientPasswordComplexity: "password doesn't match regex",
		Err400_UserWithEmailExists:           "user with such email exists",
		Err400_MalformedJSON:                 "malformed JSON request",
		Err400_InvalidRequestBody:            "invalid request body",
	},
	// 401
	http.StatusUnauthorized: {
		Err401_InvalidCredentials:         "invalid credentials provided",
		Err401_AuthorizationHeaderMissing: "authorization header missing",
		Err401_AuthorizationHeaderInvalid: "authorization header is invalid",
		Err401_AuthorizationExpired:       "session expired",
		Err401_InvalidRefreshToken:        "invalid refresh token",
		Err401_UserNotFound:               "no user found",
	},
	// 403
	http.StatusForbidden: {
		Err403_CannotToDelete:  "insufficient permissions for deletion",
		Err403_CannotEditPhone: "insufficient permissions for phone number edit",
	},
	// 404
	http.StatusNotFound: {
		Err404_PlayerStatsNotFound: "player stats not Available",
		Err404_UserOrPhoneNotFound: "user or phone number not found",
		Err404_AccountNotFound:     "account already deleted or does not exist",
	},
	// 429
	http.StatusTooManyRequests: {
		Err429_EditRequestTimedOut: "edit requests reached limit",
	},
	// 500
	http.StatusInternalServerError: {
		Err500_UnableToDelete:        "unexpected issue during deletion",
		Err500_UnableToEditPhone:     "unexpected issue during phone number edit",
		Err500_UnableToRegister:      "unexpected issue during registration",
		Err500_UnableToGenerateToken: "unable to generate JWT token",
		Err500_UnknownError:          "internal server error",
	},
	// 503
	http.StatusServiceUnavailable: {
		Err503_DataBaseOnDelete:    "Service Unavailable: Database issue during deletion",
		Err503_DataBaseOnPhoneEdit: "Service Unavailable: Database issue during phone number edit",
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
*/