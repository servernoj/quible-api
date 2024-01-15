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
	Err400_MissingRequiredQueryParam
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
	Err424_ScheduleSeason
	Err424_DailySchedule
	Err424_TeamInfo
	Err424_TeamStats
	Err424_PlayerInfo
	Err424_PlayerStats
	Err424_Injuries
	Err424_LiveFeed
	Err424_BasketAPIGetGames
)
const (
	Err500_UnknownError ErrorCode = Err500_Shift + iota + 1
)

// TODO: Complete the mapping
var ErrorMap = c.ErrorMap[ErrorCode]{
	// 400
	http.StatusBadRequest: {
		Err400_UnknownError:              "unknown error",
		Err400_MalformedJSON:             "malformed JSON request",
		Err400_InvalidRequestBody:        "invalid request body",
		Err400_MissingRequiredQueryParam: "missing required query param",
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
		Err424_UnknownError:      "unknown error",
		Err424_ScheduleSeason:    "unexpected problem with /schedule-season RSC API",
		Err424_DailySchedule:     "unexpected problem with /schedule RSC API",
		Err424_TeamInfo:          "unexpected problem with /team-info RSC API",
		Err424_TeamStats:         "unexpected problem with /team-stats RSC API",
		Err424_PlayerInfo:        "unexpected problem with /player-info RSC API",
		Err424_PlayerStats:       "unexpected problem with /player-stats RSC API",
		Err424_Injuries:          "unexpected problem with /injuries RSC API",
		Err424_LiveFeed:          "unexpected problem with /live RSC API",
		Err424_BasketAPIGetGames: "unexpected problem with /matches from BasketAPI",
	},
	// 500
	http.StatusInternalServerError: {
		Err500_UnknownError: "internal server error",
	},
}
