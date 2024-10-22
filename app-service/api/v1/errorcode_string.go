// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package v1

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Err400_UnknownError-4002001]
	_ = x[Err400_MalformedJSON-4002002]
	_ = x[Err400_InvalidRequest-4002003]
	_ = x[Err400_MissingRequiredQueryParam-4002004]
	_ = x[Err400_ChatGroupExists-4002005]
	_ = x[Err400_ChannelExists-4002006]
	_ = x[Err400_ChatGroupIsPrivate-4002007]
	_ = x[Err400_ChatGroupIsPublic-4002008]
	_ = x[Err400_ChatGroupIsSelfOwned-4002009]
	_ = x[Err400_ChatChannelAlreadyJoined-4002010]
	_ = x[Err400_EmailNotFound-4002011]
	_ = x[Err400_InvalidOrMalformedToken-4002012]
	_ = x[Err400_ChatChannelInviteeNotUser-4002013]
	_ = x[Err400_ChatChannelInviteeOwnsChatGroup-4002014]
	_ = x[Err400_OnlyForChatGroups-4002015]
	_ = x[Err400_OnlyForChatChannels-4002016]
	_ = x[Err401_UnknownError-4012001]
	_ = x[Err401_UserIdNotFound-4012002]
	_ = x[Err401_UserNotFound-4012003]
	_ = x[Err401_AuthServiceError-4012004]
	_ = x[Err401_InvalidAccessToken-4012005]
	_ = x[Err404_UnknownError-4042001]
	_ = x[Err404_ChatGroupNotFound-4042002]
	_ = x[Err404_ChatChannelNotFound-4042003]
	_ = x[Err404_ChatRecordNotFound-4042004]
	_ = x[Err417_UnknownError-4172001]
	_ = x[Err417_InvalidToken-4172002]
	_ = x[Err424_UnknownError-4242001]
	_ = x[Err424_ScheduleSeason-4242002]
	_ = x[Err424_DailySchedule-4242003]
	_ = x[Err424_TeamInfo-4242004]
	_ = x[Err424_TeamStats-4242005]
	_ = x[Err424_PlayerInfo-4242006]
	_ = x[Err424_PlayerStats-4242007]
	_ = x[Err424_Injuries-4242008]
	_ = x[Err424_LiveFeed-4242009]
	_ = x[Err424_BasketAPIListGames-4242010]
	_ = x[Err424_BasketAPIGetGame-4242011]
	_ = x[Err424_UnableToSendEmail-4242012]
	_ = x[Err500_UnknownError-5002001]
	_ = x[Err500_UnknownHumaError-5002002]
	_ = x[Err500_UnableCreateChatUser-5002003]
	_ = x[Err500_UnableUpdateChatUser-5002004]
	_ = x[Err500_UnableUpdateChatRecord-5002005]
}

const (
	_ErrorCode_name_0 = "Err400_UnknownErrorErr400_MalformedJSONErr400_InvalidRequestErr400_MissingRequiredQueryParamErr400_ChatGroupExistsErr400_ChannelExistsErr400_ChatGroupIsPrivateErr400_ChatGroupIsPublicErr400_ChatGroupIsSelfOwnedErr400_ChatChannelAlreadyJoinedErr400_EmailNotFoundErr400_InvalidOrMalformedTokenErr400_ChatChannelInviteeNotUserErr400_ChatChannelInviteeOwnsChatGroupErr400_OnlyForChatGroupsErr400_OnlyForChatChannels"
	_ErrorCode_name_1 = "Err401_UnknownErrorErr401_UserIdNotFoundErr401_UserNotFoundErr401_AuthServiceErrorErr401_InvalidAccessToken"
	_ErrorCode_name_2 = "Err404_UnknownErrorErr404_ChatGroupNotFoundErr404_ChatChannelNotFoundErr404_ChatRecordNotFound"
	_ErrorCode_name_3 = "Err417_UnknownErrorErr417_InvalidToken"
	_ErrorCode_name_4 = "Err424_UnknownErrorErr424_ScheduleSeasonErr424_DailyScheduleErr424_TeamInfoErr424_TeamStatsErr424_PlayerInfoErr424_PlayerStatsErr424_InjuriesErr424_LiveFeedErr424_BasketAPIListGamesErr424_BasketAPIGetGameErr424_UnableToSendEmail"
	_ErrorCode_name_5 = "Err500_UnknownErrorErr500_UnknownHumaErrorErr500_UnableCreateChatUserErr500_UnableUpdateChatUserErr500_UnableUpdateChatRecord"
)

var (
	_ErrorCode_index_0 = [...]uint16{0, 19, 39, 60, 92, 114, 134, 159, 183, 210, 241, 261, 291, 323, 361, 385, 411}
	_ErrorCode_index_1 = [...]uint8{0, 19, 40, 59, 82, 107}
	_ErrorCode_index_2 = [...]uint8{0, 19, 43, 69, 94}
	_ErrorCode_index_3 = [...]uint8{0, 19, 38}
	_ErrorCode_index_4 = [...]uint8{0, 19, 40, 60, 75, 91, 108, 126, 141, 156, 181, 204, 228}
	_ErrorCode_index_5 = [...]uint8{0, 19, 42, 69, 96, 125}
)

func (i ErrorCode) String() string {
	switch {
	case 4002001 <= i && i <= 4002016:
		i -= 4002001
		return _ErrorCode_name_0[_ErrorCode_index_0[i]:_ErrorCode_index_0[i+1]]
	case 4012001 <= i && i <= 4012005:
		i -= 4012001
		return _ErrorCode_name_1[_ErrorCode_index_1[i]:_ErrorCode_index_1[i+1]]
	case 4042001 <= i && i <= 4042004:
		i -= 4042001
		return _ErrorCode_name_2[_ErrorCode_index_2[i]:_ErrorCode_index_2[i+1]]
	case 4172001 <= i && i <= 4172002:
		i -= 4172001
		return _ErrorCode_name_3[_ErrorCode_index_3[i]:_ErrorCode_index_3[i+1]]
	case 4242001 <= i && i <= 4242012:
		i -= 4242001
		return _ErrorCode_name_4[_ErrorCode_index_4[i]:_ErrorCode_index_4[i+1]]
	case 5002001 <= i && i <= 5002005:
		i -= 5002001
		return _ErrorCode_name_5[_ErrorCode_index_5[i]:_ErrorCode_index_5[i+1]]
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
