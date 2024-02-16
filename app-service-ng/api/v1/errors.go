package v1

import (
	"net/http"

	libAPI "github.com/quible-io/quible-api/lib/api"
)

const ErrServiceId = 2000

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	Err400_Shift = libAPI.ErrStatusGain*http.StatusBadRequest + ErrServiceId
	Err401_Shift = libAPI.ErrStatusGain*http.StatusUnauthorized + ErrServiceId
	Err404_Shift = libAPI.ErrStatusGain*http.StatusNotFound + ErrServiceId
	Err417_Shift = libAPI.ErrStatusGain*http.StatusExpectationFailed + ErrServiceId
	Err424_Shift = libAPI.ErrStatusGain*http.StatusFailedDependency + ErrServiceId
	Err500_Shift = libAPI.ErrStatusGain*http.StatusInternalServerError + ErrServiceId
)

const (
	Err400_UnknownError ErrorCode = Err400_Shift + iota + 1
	Err400_MalformedJSON
	Err400_InvalidRequestBody
	Err400_MissingRequiredQueryParam
	Err400_ChatGroupExists
	Err400_ChannelExists
	Err400_ChatGroupIsPrivate
	Err400_ChatGroupIsPublic
	Err400_ChatGroupIsSelfOwned
	Err400_ChannelAlreadyJoined
	Err400_EmailNotFound
	Err400_InvalidOrMalformedToken
)
const (
	Err401_UnknownError ErrorCode = Err401_Shift + iota + 1
	Err401_UserIdNotFound
	Err401_UserNotFound
	Err401_InvalidAccessToken
)
const (
	Err404_UnknownError ErrorCode = Err404_Shift + iota + 1
	Err404_ChatGroupNotFound
	Err404_ChannelNotFound
)
const (
	Err417_UnknownError ErrorCode = Err417_Shift + iota + 1
	Err417_InvalidToken
)
const (
	Err424_UnknownError   ErrorCode = Err424_Shift + iota + 1
	Err424_ScheduleSeason           // unused
	Err424_DailySchedule            // unused
	Err424_TeamInfo                 // unused
	Err424_TeamStats                // unused
	Err424_PlayerInfo               // unused
	Err424_PlayerStats              // unused
	Err424_Injuries                 // unused
	Err424_LiveFeed                 // unused
	Err424_BasketAPIGetGames
	Err424_BasketAPIGetGameDetails
	Err424_UnableToSendEmail
)
const (
	Err500_UnknownError ErrorCode = Err500_Shift + iota + 1
	Err500_UnknownHumaError
)

var ErrorMap = libAPI.ErrorMap[ErrorCode]{
	// 400
	Err400_UnknownError:              "unknown error",
	Err400_MalformedJSON:             "malformed JSON request",
	Err400_InvalidRequestBody:        "invalid request body",
	Err400_MissingRequiredQueryParam: "missing/invalid required query param",
	Err400_ChatGroupExists:           "chat group with given name or title exists",
	Err400_ChannelExists:             "channel with this name in the same chat group already exists",
	Err400_ChatGroupIsPrivate:        "chat group holding the channel is private",
	Err400_ChatGroupIsSelfOwned:      "chat group holding the channel is self-owned",
	Err400_ChannelAlreadyJoined:      "channel already joined",
	Err400_ChatGroupIsPublic:         "chat group holding the channel is public",
	Err400_EmailNotFound:             "email not found",
	Err400_InvalidOrMalformedToken:   "activation token is missing or malformed",
	// 401
	Err401_UnknownError:       "unknown error",
	Err401_UserIdNotFound:     "userId not present",
	Err401_UserNotFound:       "user not found",
	Err401_InvalidAccessToken: "invalid or missing access token",
	// 404
	Err404_UnknownError:      "unknown error",
	Err404_ChatGroupNotFound: "chat group not found",
	Err404_ChannelNotFound:   "channel not found",
	// 417
	Err417_UnknownError: "unknown error",
	Err417_InvalidToken: "invalid (possibly expired) token",
	// 424
	Err424_UnknownError:            "unknown error",
	Err424_BasketAPIGetGames:       "unexpected problem with MatchSchedules from BasketAPI",
	Err424_BasketAPIGetGameDetails: "unexpected problem with (Match|MatchStatistics|MatchLineups) API from BasketAPI",
	Err424_UnableToSendEmail:       "unable to send email",
	// 500
	Err500_UnknownError:     "internal server error",
	Err500_UnknownHumaError: "unidentified upstream Huma error",
}
