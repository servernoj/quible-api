package v1

import (
	"net/http"

	"github.com/quible-io/quible-api/auth-service/api"
)

const ErrServiceId = 1000

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	Err207_Shift = api.ErrStatusGain*http.StatusMultiStatus + ErrServiceId
	Err400_Shift = api.ErrStatusGain*http.StatusBadRequest + ErrServiceId
	Err401_Shift = api.ErrStatusGain*http.StatusUnauthorized + ErrServiceId
	Err403_Shift = api.ErrStatusGain*http.StatusForbidden + ErrServiceId
	Err404_Shift = api.ErrStatusGain*http.StatusNotFound + ErrServiceId
	Err417_Shift = api.ErrStatusGain*http.StatusExpectationFailed + ErrServiceId
	Err422_Shift = api.ErrStatusGain*http.StatusUnprocessableEntity + ErrServiceId
	Err424_Shift = api.ErrStatusGain*http.StatusFailedDependency + ErrServiceId
	Err429_Shift = api.ErrStatusGain*http.StatusTooManyRequests + ErrServiceId
	Err500_Shift = api.ErrStatusGain*http.StatusInternalServerError + ErrServiceId
	Err503_Shift = api.ErrStatusGain*http.StatusServiceUnavailable + ErrServiceId
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
	Err400_InsufficientPasswordComplexity
	Err400_MalformedJSON
	Err400_InvalidRequestBody
	Err400_FileTooLarge
	Err400_InvalidClientId
	Err400_UserWithEmailOrUsernameExists
	Err400_InvalidOrMalformedToken
	Err400_ImageDataNotPresent
	Err400_UnsatisfactoryPassword
	Err400_UnsatisfactoryConfirmPassword
	Err400_UserWithEmailExists
)
const (
	Err401_InvalidCredentials ErrorCode = Err401_Shift + iota + 1
	Err401_AuthorizationHeaderMissing
	Err401_AuthorizationHeaderInvalid
	Err401_AuthorizationExpired
	Err401_InvalidRefreshToken
	Err401_UserNotFound
	Err401_UserNotActivated
	Err401_InvalidAccessToken
	Err401_InvalidActivationToken
	Err401_InvalidPasswordResetToken
)
const (
	Err403_CannotToDelete ErrorCode = Err403_Shift + iota + 1
	Err403_CannotEditPhone
)
const (
	Err404_PlayerStatsNotFound ErrorCode = Err404_Shift + iota + 1
	Err404_UserOrPhoneNotFound
	Err404_AccountNotFound
	Err404_UserNotFound
	Err404_UserHasNoImage
)
const (
	Err417_UnknownError ErrorCode = Err417_Shift + iota + 1
	Err417_InvalidToken
	Err417_UnableToAssociateUser
)
const (
	Err422_UnknownError ErrorCode = Err422_Shift + iota + 1
)
const (
	Err424_UnknownError ErrorCode = Err424_Shift + iota + 1
	Err424_UnableToSendEmail
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
	Err500_UnableToResetPassword
	Err500_UnableToActivateUser
	Err500_UnableToUpdateUser
	Err500_UnknownHumaError
	Err500_UnableToRetrieveProfileImage
	Err500_UnableToStoreImage
)
const (
	Err503_DataBaseOnDelete ErrorCode = Err503_Shift + iota + 1
	Err503_DataBaseOnPhoneEdit
)

var ErrorMap = api.ErrorMap[ErrorCode]{
	// -- 400
	Err400_EmailNotRegistered:             "email is not registered",
	Err400_InvalidEmailFormat:             "invalid email address format",
	Err400_InvalidPhoneFormat:             "invalid phone number format",
	Err400_MalformedJSON:                  "malformed JSON request",
	Err400_InvalidRequestBody:             "invalid request body",
	Err400_FileTooLarge:                   "invalid file size",
	Err400_InvalidClientId:                "unexpected clientId",
	Err400_UserWithEmailOrUsernameExists:  "activated user with such username or email exists",
	Err400_InvalidOrMalformedToken:        "activation token is missing or malformed",
	Err400_InsufficientPasswordComplexity: "insufficient password complexity",
	Err400_ImageDataNotPresent:            "image data not present in multipart request body under key `image`",
	Err400_UnsatisfactoryPassword:         "unsatisfactory value of the password field",
	Err400_UnsatisfactoryConfirmPassword:  "unsatisfactory value of the confirmPassword field",
	Err400_UserWithEmailExists:            "user with provided email already exists",

	// -- 401
	Err401_InvalidCredentials:         "invalid credentials provided",
	Err401_AuthorizationHeaderMissing: "authorization header missing",
	Err401_AuthorizationHeaderInvalid: "authorization header is invalid",
	Err401_AuthorizationExpired:       "session expired",
	Err401_InvalidRefreshToken:        "invalid refresh token",
	Err401_UserNotFound:               "user not found",
	Err401_UserNotActivated:           "user is not activated",
	Err401_InvalidAccessToken:         "invalid or missing access token",
	Err401_InvalidActivationToken:     "invalid or missing activation token",
	Err401_InvalidPasswordResetToken:  "invalid or missing password reset token",
	// -- 404
	Err404_UserNotFound:   "user not found",
	Err404_UserHasNoImage: "user has no profile image",
	// -- 417
	Err417_InvalidToken:          "invalid (possibly expired) token",
	Err417_UnableToAssociateUser: "unable to associate user with the token",
	// -- 424
	Err424_UnableToSendEmail: "unable to send email",
	// -- 500
	Err500_UnableToRegister:             "unexpected issue during registration",
	Err500_UnableToGenerateToken:        "unable to generate JWT token",
	Err500_UnknownError:                 "internal server error",
	Err500_UnableToActivateUser:         "unable to activate user",
	Err500_UnableToResetPassword:        "unable to reset password",
	Err500_UnableToUpdateUser:           "unable to update user record",
	Err500_UnknownHumaError:             "unidentified upstream Huma error",
	Err500_UnableToRetrieveProfileImage: "unable to retrieve profile image",
	Err500_UnableToStoreImage:           "unable to store uploaded profile image",
}
