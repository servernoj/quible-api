// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package controller

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Err207_SomeDataUndeleted-2071001]
	_ = x[Err400_EmailNotRegistered-4001001]
	_ = x[Err400_InvalidEmailFormat-4001002]
	_ = x[Err400_InvalidUsernameFormat-4001003]
	_ = x[Err400_InvalidPhoneFormat-4001004]
	_ = x[Err400_UserWithUsernameExists-4001005]
	_ = x[Err400_UserWithEmailExists-4001006]
	_ = x[Err400_IsufficientPasswordComplexity-4001007]
	_ = x[Err400_MalformedJSON-4001008]
	_ = x[Err400_InvalidRequestBody-4001009]
	_ = x[Err400_FileTooLarge-4001010]
	_ = x[Err400_InvalidClientId-4001011]
	_ = x[Err400_UserWithEmailOrUsernameExists-4001012]
	_ = x[Err400_InvalidOrMalformedToken-4001013]
	_ = x[Err400_ChatGroupExists-4001014]
	_ = x[Err400_ChannelExists-4001015]
	_ = x[Err400_ChatGroupIsPrivate-4001016]
	_ = x[Err400_ChatGroupIsSelfOwned-4001017]
	_ = x[Err400_ChannelAlreadyJoined-4001018]
	_ = x[Err401_InvalidCredentials-4011001]
	_ = x[Err401_AuthorizationHeaderMissing-4011002]
	_ = x[Err401_AuthorizationHeaderInvalid-4011003]
	_ = x[Err401_AuthorizationExpired-4011004]
	_ = x[Err401_InvalidRefreshToken-4011005]
	_ = x[Err401_UserNotFound-4011006]
	_ = x[Err401_UserNotActivated-4011007]
	_ = x[Err403_CannotToDelete-4031001]
	_ = x[Err403_CannotEditPhone-4031002]
	_ = x[Err404_PlayerStatsNotFound-4041001]
	_ = x[Err404_UserOrPhoneNotFound-4041002]
	_ = x[Err404_AccountNotFound-4041003]
	_ = x[Err404_UserNotFound-4041004]
	_ = x[Err404_UserHasNoImage-4041005]
	_ = x[Err404_ChatGroupNotFound-4041006]
	_ = x[Err404_ChannelNotFound-4041007]
	_ = x[Err417_UnknownError-4171001]
	_ = x[Err417_InvalidToken-4171002]
	_ = x[Err417_UnableToAssociateUser-4171003]
	_ = x[Err424_UnknownError-4241001]
	_ = x[Err424_UnableToSendEmail-4241002]
	_ = x[Err429_EditRequestTimedOut-4291001]
	_ = x[Err500_UnknownError-5001001]
	_ = x[Err500_UnableToDelete-5001002]
	_ = x[Err500_UnableToEditPhone-5001003]
	_ = x[Err500_UnableToRegister-5001004]
	_ = x[Err500_UnableToGenerateToken-5001005]
	_ = x[Err500_UnableToResetPassword-5001006]
	_ = x[Err500_UnableToActivateUser-5001007]
	_ = x[Err503_DataBaseOnDelete-5031001]
	_ = x[Err503_DataBaseOnPhoneEdit-5031002]
}

const (
	_ErrorCode_name_0 = "Err207_SomeDataUndeleted"
	_ErrorCode_name_1 = "Err400_EmailNotRegisteredErr400_InvalidEmailFormatErr400_InvalidUsernameFormatErr400_InvalidPhoneFormatErr400_UserWithUsernameExistsErr400_UserWithEmailExistsErr400_IsufficientPasswordComplexityErr400_MalformedJSONErr400_InvalidRequestBodyErr400_FileTooLargeErr400_InvalidClientIdErr400_UserWithEmailOrUsernameExistsErr400_InvalidOrMalformedTokenErr400_ChatGroupExistsErr400_ChannelExistsErr400_ChatGroupIsPrivateErr400_ChatGroupIsSelfOwnedErr400_ChannelAlreadyJoined"
	_ErrorCode_name_2 = "Err401_InvalidCredentialsErr401_AuthorizationHeaderMissingErr401_AuthorizationHeaderInvalidErr401_AuthorizationExpiredErr401_InvalidRefreshTokenErr401_UserNotFoundErr401_UserNotActivated"
	_ErrorCode_name_3 = "Err403_CannotToDeleteErr403_CannotEditPhone"
	_ErrorCode_name_4 = "Err404_PlayerStatsNotFoundErr404_UserOrPhoneNotFoundErr404_AccountNotFoundErr404_UserNotFoundErr404_UserHasNoImageErr404_ChatGroupNotFoundErr404_ChannelNotFound"
	_ErrorCode_name_5 = "Err417_UnknownErrorErr417_InvalidTokenErr417_UnableToAssociateUser"
	_ErrorCode_name_6 = "Err424_UnknownErrorErr424_UnableToSendEmail"
	_ErrorCode_name_7 = "Err429_EditRequestTimedOut"
	_ErrorCode_name_8 = "Err500_UnknownErrorErr500_UnableToDeleteErr500_UnableToEditPhoneErr500_UnableToRegisterErr500_UnableToGenerateTokenErr500_UnableToResetPasswordErr500_UnableToActivateUser"
	_ErrorCode_name_9 = "Err503_DataBaseOnDeleteErr503_DataBaseOnPhoneEdit"
)

var (
	_ErrorCode_index_1 = [...]uint16{0, 25, 50, 78, 103, 132, 158, 194, 214, 239, 258, 280, 316, 346, 368, 388, 413, 440, 467}
	_ErrorCode_index_2 = [...]uint8{0, 25, 58, 91, 118, 144, 163, 186}
	_ErrorCode_index_3 = [...]uint8{0, 21, 43}
	_ErrorCode_index_4 = [...]uint8{0, 26, 52, 74, 93, 114, 138, 160}
	_ErrorCode_index_5 = [...]uint8{0, 19, 38, 66}
	_ErrorCode_index_6 = [...]uint8{0, 19, 43}
	_ErrorCode_index_8 = [...]uint8{0, 19, 40, 64, 87, 115, 143, 170}
	_ErrorCode_index_9 = [...]uint8{0, 23, 49}
)

func (i ErrorCode) String() string {
	switch {
	case i == 2071001:
		return _ErrorCode_name_0
	case 4001001 <= i && i <= 4001018:
		i -= 4001001
		return _ErrorCode_name_1[_ErrorCode_index_1[i]:_ErrorCode_index_1[i+1]]
	case 4011001 <= i && i <= 4011007:
		i -= 4011001
		return _ErrorCode_name_2[_ErrorCode_index_2[i]:_ErrorCode_index_2[i+1]]
	case 4031001 <= i && i <= 4031002:
		i -= 4031001
		return _ErrorCode_name_3[_ErrorCode_index_3[i]:_ErrorCode_index_3[i+1]]
	case 4041001 <= i && i <= 4041007:
		i -= 4041001
		return _ErrorCode_name_4[_ErrorCode_index_4[i]:_ErrorCode_index_4[i+1]]
	case 4171001 <= i && i <= 4171003:
		i -= 4171001
		return _ErrorCode_name_5[_ErrorCode_index_5[i]:_ErrorCode_index_5[i+1]]
	case 4241001 <= i && i <= 4241002:
		i -= 4241001
		return _ErrorCode_name_6[_ErrorCode_index_6[i]:_ErrorCode_index_6[i+1]]
	case i == 4291001:
		return _ErrorCode_name_7
	case 5001001 <= i && i <= 5001007:
		i -= 5001001
		return _ErrorCode_name_8[_ErrorCode_index_8[i]:_ErrorCode_index_8[i+1]]
	case 5031001 <= i && i <= 5031002:
		i -= 5031001
		return _ErrorCode_name_9[_ErrorCode_index_9[i]:_ErrorCode_index_9[i+1]]
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
