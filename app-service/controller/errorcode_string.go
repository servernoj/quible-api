// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package controller

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Err400_UnknownError-4002001]
	_ = x[Err400_MalformedJSON-4002002]
	_ = x[Err400_InvalidRequestBody-4002003]
	_ = x[Err401_UnknownError-4012001]
	_ = x[Err401_UserIdNotFound-4012002]
	_ = x[Err401_UserNotFound-4012003]
	_ = x[Err404_UnknownError-4042001]
	_ = x[Err424_UnknownError-4242001]
	_ = x[Err424_ScheduleSeason-4242002]
	_ = x[Err424_DailySchedule-4242003]
	_ = x[Err424_TeamInfo-4242004]
	_ = x[Err424_TeamStats-4242005]
	_ = x[Err424_PlayerInfo-4242006]
	_ = x[Err424_PlayerStats-4242007]
	_ = x[Err424_Injuries-4242008]
	_ = x[Err424_LiveFeed-4242009]
	_ = x[Err500_UnknownError-5002001]
}

const (
	_ErrorCode_name_0 = "Err400_UnknownErrorErr400_MalformedJSONErr400_InvalidRequestBody"
	_ErrorCode_name_1 = "Err401_UnknownErrorErr401_UserIdNotFoundErr401_UserNotFound"
	_ErrorCode_name_2 = "Err404_UnknownError"
	_ErrorCode_name_3 = "Err424_UnknownErrorErr424_ScheduleSeasonErr424_DailyScheduleErr424_TeamInfoErr424_TeamStatsErr424_PlayerInfoErr424_PlayerStatsErr424_InjuriesErr424_LiveFeed"
	_ErrorCode_name_4 = "Err500_UnknownError"
)

var (
	_ErrorCode_index_0 = [...]uint8{0, 19, 39, 64}
	_ErrorCode_index_1 = [...]uint8{0, 19, 40, 59}
	_ErrorCode_index_3 = [...]uint8{0, 19, 40, 60, 75, 91, 108, 126, 141, 156}
)

func (i ErrorCode) String() string {
	switch {
	case 4002001 <= i && i <= 4002003:
		i -= 4002001
		return _ErrorCode_name_0[_ErrorCode_index_0[i]:_ErrorCode_index_0[i+1]]
	case 4012001 <= i && i <= 4012003:
		i -= 4012001
		return _ErrorCode_name_1[_ErrorCode_index_1[i]:_ErrorCode_index_1[i+1]]
	case i == 4042001:
		return _ErrorCode_name_2
	case 4242001 <= i && i <= 4242009:
		i -= 4242001
		return _ErrorCode_name_3[_ErrorCode_index_3[i]:_ErrorCode_index_3[i+1]]
	case i == 5002001:
		return _ErrorCode_name_4
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}