package v1

import (
	"net/http"

	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/misc"
)

const ErrServiceId = 2000

//go:generate stringer -type=ErrorCode

type ErrorCode int

func (ec ErrorCode) Ptr() *int {
	return misc.Of(int(ec))
}

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
	Err400_InvalidRequest
	Err400_MissingRequiredQueryParam
	Err400_ChatGroupExists
	Err400_ChannelExists
	Err400_ChatGroupIsPrivate
	Err400_ChatGroupIsPublic
	Err400_ChatGroupIsSelfOwned
	Err400_ChatChannelAlreadyJoined
	Err400_EmailNotFound
	Err400_InvalidOrMalformedToken
	Err400_ChatChannelInviteeNotUser
	Err400_ChatChannelInviteeOwnsChatGroup
	Err400_OnlyForChatGroups
	Err400_OnlyForChatChannels
)
const (
	Err401_UnknownError ErrorCode = Err401_Shift + iota + 1
	Err401_UserIdNotFound
	Err401_UserNotFound
	Err401_AuthServiceError
	Err401_InvalidAccessToken
)
const (
	Err404_UnknownError ErrorCode = Err404_Shift + iota + 1
	Err404_ChatGroupNotFound
	Err404_ChatChannelNotFound
	Err404_ChatRecordNotFound
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
	Err424_BasketAPIListGames
	Err424_BasketAPIGetGame
	Err424_UnableToSendEmail
)
const (
	Err500_UnknownError ErrorCode = Err500_Shift + iota + 1
	Err500_UnknownHumaError
	Err500_UnableCreateChatUser
	Err500_UnableUpdateChatUser
	Err500_UnableUpdateChatRecord
)

var ErrorMap = libAPI.ErrorMap[ErrorCode]{
	// 400
	Err400_UnknownError:                    "unknown error",
	Err400_MalformedJSON:                   "malformed JSON request",
	Err400_InvalidRequest:                  "invalid request",
	Err400_MissingRequiredQueryParam:       "missing/invalid required query param",
	Err400_ChatGroupExists:                 "chat group with given name or title exists",
	Err400_ChannelExists:                   "chat channel with this name in the same chat group already exists",
	Err400_ChatGroupIsPrivate:              "chat group holding the channel is private",
	Err400_ChatGroupIsSelfOwned:            "chat group holding the channel is self-owned",
	Err400_ChatChannelAlreadyJoined:        "chat channel already joined",
	Err400_ChatGroupIsPublic:               "chat group holding the channel is public",
	Err400_EmailNotFound:                   "email not found",
	Err400_InvalidOrMalformedToken:         "activation token is missing or malformed",
	Err400_ChatChannelInviteeNotUser:       "chat channel invitee doesn't have an account",
	Err400_ChatChannelInviteeOwnsChatGroup: "chat group owner cannot be an invitee",
	Err400_OnlyForChatGroups:               "update allowed only for chat groups",
	Err400_OnlyForChatChannels:             "update allowed only for chat channels",
	// 401
	Err401_UnknownError:       "unknown error",
	Err401_UserIdNotFound:     "userId not present",
	Err401_UserNotFound:       "user not found",
	Err401_InvalidAccessToken: "invalid or missing access token",
	Err401_AuthServiceError:   "unexpected auth-service failure",
	// 404
	Err404_UnknownError:        "unknown error",
	Err404_ChatGroupNotFound:   "chat group not found",
	Err404_ChatChannelNotFound: "chat channel not found",
	Err404_ChatRecordNotFound:  "chat record (group/channel) not found",
	// 417
	Err417_UnknownError: "unknown error",
	Err417_InvalidToken: "invalid (possibly expired) token",
	// 424
	Err424_UnknownError:       "unknown error",
	Err424_BasketAPIListGames: "unexpected problem with MatchSchedules from BasketAPI",
	Err424_BasketAPIGetGame:   "unexpected problem with (Match|MatchStatistics|MatchLineups) API from BasketAPI",
	Err424_UnableToSendEmail:  "unable to send email",
	// 500
	Err500_UnknownError:           "internal server error",
	Err500_UnknownHumaError:       "unidentified upstream Huma error",
	Err500_UnableCreateChatUser:   "unable to create chat to user association",
	Err500_UnableUpdateChatUser:   "unable to update chat to user association",
	Err500_UnableUpdateChatRecord: "unable to update chat record (group/channel)",
}
