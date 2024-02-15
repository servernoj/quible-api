// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package v1

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
	_ = x[Err400_InsufficientPasswordComplexity-4001006]
	_ = x[Err400_MalformedJSON-4001007]
	_ = x[Err400_InvalidRequestBody-4001008]
	_ = x[Err400_FileTooLarge-4001009]
	_ = x[Err400_InvalidClientId-4001010]
	_ = x[Err400_UserWithEmailOrUsernameExists-4001011]
	_ = x[Err400_InvalidOrMalformedToken-4001012]
	_ = x[Err400_ImageDataNotPresent-4001013]
	_ = x[Err400_UnsatisfactoryPassword-4001014]
	_ = x[Err400_UnsatisfactoryConfirmPassword-4001015]
	_ = x[Err400_UserWithEmailExists-4001016]
	_ = x[Err401_InvalidCredentials-4011001]
	_ = x[Err401_AuthorizationHeaderMissing-4011002]
	_ = x[Err401_AuthorizationHeaderInvalid-4011003]
	_ = x[Err401_AuthorizationExpired-4011004]
	_ = x[Err401_InvalidRefreshToken-4011005]
	_ = x[Err401_UserNotFound-4011006]
	_ = x[Err401_UserNotActivated-4011007]
	_ = x[Err401_InvalidAccessToken-4011008]
	_ = x[Err401_InvalidActivationToken-4011009]
	_ = x[Err401_InvalidPasswordResetToken-4011010]
	_ = x[Err403_CannotToDelete-4031001]
	_ = x[Err403_CannotEditPhone-4031002]
	_ = x[Err404_PlayerStatsNotFound-4041001]
	_ = x[Err404_UserOrPhoneNotFound-4041002]
	_ = x[Err404_AccountNotFound-4041003]
	_ = x[Err404_UserNotFound-4041004]
	_ = x[Err404_UserHasNoImage-4041005]
	_ = x[Err417_UnknownError-4171001]
	_ = x[Err417_InvalidToken-4171002]
	_ = x[Err417_UnableToAssociateUser-4171003]
	_ = x[Err422_UnknownError-4221001]
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
	_ = x[Err500_UnableToUpdateUser-5001008]
	_ = x[Err500_UnknownHumaError-5001009]
	_ = x[Err500_UnableToRetrieveProfileImage-5001010]
	_ = x[Err500_UnableToStoreImage-5001011]
	_ = x[Err503_DataBaseOnDelete-5031001]
	_ = x[Err503_DataBaseOnPhoneEdit-5031002]
}

const _ErrorCode_name = "Err207_SomeDataUndeletedErr400_EmailNotRegisteredErr400_InvalidEmailFormatErr400_InvalidUsernameFormatErr400_InvalidPhoneFormatErr400_UserWithUsernameExistsErr400_InsufficientPasswordComplexityErr400_MalformedJSONErr400_InvalidRequestBodyErr400_FileTooLargeErr400_InvalidClientIdErr400_UserWithEmailOrUsernameExistsErr400_InvalidOrMalformedTokenErr400_ImageDataNotPresentErr400_UnsatisfactoryPasswordErr400_UnsatisfactoryConfirmPasswordErr400_UserWithEmailExistsErr401_InvalidCredentialsErr401_AuthorizationHeaderMissingErr401_AuthorizationHeaderInvalidErr401_AuthorizationExpiredErr401_InvalidRefreshTokenErr401_UserNotFoundErr401_UserNotActivatedErr401_InvalidAccessTokenErr401_InvalidActivationTokenErr401_InvalidPasswordResetTokenErr403_CannotToDeleteErr403_CannotEditPhoneErr404_PlayerStatsNotFoundErr404_UserOrPhoneNotFoundErr404_AccountNotFoundErr404_UserNotFoundErr404_UserHasNoImageErr417_UnknownErrorErr417_InvalidTokenErr417_UnableToAssociateUserErr422_UnknownErrorErr424_UnknownErrorErr424_UnableToSendEmailErr429_EditRequestTimedOutErr500_UnknownErrorErr500_UnableToDeleteErr500_UnableToEditPhoneErr500_UnableToRegisterErr500_UnableToGenerateTokenErr500_UnableToResetPasswordErr500_UnableToActivateUserErr500_UnableToUpdateUserErr500_UnknownHumaErrorErr500_UnableToRetrieveProfileImageErr500_UnableToStoreImageErr503_DataBaseOnDeleteErr503_DataBaseOnPhoneEdit"

var _ErrorCode_map = map[ErrorCode]string{
	2071001: _ErrorCode_name[0:24],
	4001001: _ErrorCode_name[24:49],
	4001002: _ErrorCode_name[49:74],
	4001003: _ErrorCode_name[74:102],
	4001004: _ErrorCode_name[102:127],
	4001005: _ErrorCode_name[127:156],
	4001006: _ErrorCode_name[156:193],
	4001007: _ErrorCode_name[193:213],
	4001008: _ErrorCode_name[213:238],
	4001009: _ErrorCode_name[238:257],
	4001010: _ErrorCode_name[257:279],
	4001011: _ErrorCode_name[279:315],
	4001012: _ErrorCode_name[315:345],
	4001013: _ErrorCode_name[345:371],
	4001014: _ErrorCode_name[371:400],
	4001015: _ErrorCode_name[400:436],
	4001016: _ErrorCode_name[436:462],
	4011001: _ErrorCode_name[462:487],
	4011002: _ErrorCode_name[487:520],
	4011003: _ErrorCode_name[520:553],
	4011004: _ErrorCode_name[553:580],
	4011005: _ErrorCode_name[580:606],
	4011006: _ErrorCode_name[606:625],
	4011007: _ErrorCode_name[625:648],
	4011008: _ErrorCode_name[648:673],
	4011009: _ErrorCode_name[673:702],
	4011010: _ErrorCode_name[702:734],
	4031001: _ErrorCode_name[734:755],
	4031002: _ErrorCode_name[755:777],
	4041001: _ErrorCode_name[777:803],
	4041002: _ErrorCode_name[803:829],
	4041003: _ErrorCode_name[829:851],
	4041004: _ErrorCode_name[851:870],
	4041005: _ErrorCode_name[870:891],
	4171001: _ErrorCode_name[891:910],
	4171002: _ErrorCode_name[910:929],
	4171003: _ErrorCode_name[929:957],
	4221001: _ErrorCode_name[957:976],
	4241001: _ErrorCode_name[976:995],
	4241002: _ErrorCode_name[995:1019],
	4291001: _ErrorCode_name[1019:1045],
	5001001: _ErrorCode_name[1045:1064],
	5001002: _ErrorCode_name[1064:1085],
	5001003: _ErrorCode_name[1085:1109],
	5001004: _ErrorCode_name[1109:1132],
	5001005: _ErrorCode_name[1132:1160],
	5001006: _ErrorCode_name[1160:1188],
	5001007: _ErrorCode_name[1188:1215],
	5001008: _ErrorCode_name[1215:1240],
	5001009: _ErrorCode_name[1240:1263],
	5001010: _ErrorCode_name[1263:1298],
	5001011: _ErrorCode_name[1298:1323],
	5031001: _ErrorCode_name[1323:1346],
	5031002: _ErrorCode_name[1346:1372],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
