// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package controller

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Err207_SomeDataUndeleted-2070001]
	_ = x[Err400_EmailNotRegistered-4000001]
	_ = x[Err400_InvalidEmailFormat-4000002]
	_ = x[Err400_InvalidUsernameFormat-4000003]
	_ = x[Err400_InvalidPhoneFormat-4000004]
	_ = x[Err400_UserWithUsernameExists-4000005]
	_ = x[Err400_UserWithEmailExists-4000006]
	_ = x[Err400_IsufficientPasswordComplexity-4000007]
	_ = x[Err400_MalformedJSON-4000008]
	_ = x[Err400_InvalidRequestBody-4000009]
	_ = x[Err401_InvalidCredentials-4010001]
	_ = x[Err401_AuthorizationHeaderMissing-4010002]
	_ = x[Err401_AuthorizationHeaderInvalid-4010003]
	_ = x[Err401_UserNotFound-4010004]
	_ = x[Err403_CannotToDelete-4030001]
	_ = x[Err403_CannotEditPhone-4030002]
	_ = x[Err404_PlayerStatsNotFound-4040001]
	_ = x[Err404_UserOrPhoneNotFound-4040002]
	_ = x[Err404_AccountNotFound-4040003]
	_ = x[Err429_EditRequestTimedOut-4290001]
	_ = x[Err500_UnableToDelete-5000001]
	_ = x[Err500_UnableToEditPhone-5000002]
	_ = x[Err500_UnableToRegister-5000003]
	_ = x[Err500_UnknownError-5000004]
	_ = x[Err503_DataBaseOnDelete-5030001]
	_ = x[Err503_DataBaseOnPhoneEdit-5030002]
}

const (
	_ErrorCode_name_0 = "Err207_SomeDataUndeleted"
	_ErrorCode_name_1 = "Err400_EmailNotRegisteredErr400_InvalidEmailFormatErr400_InvalidUsernameFormatErr400_InvalidPhoneFormatErr400_UserWithUsernameExistsErr400_UserWithEmailExistsErr400_IsufficientPasswordComplexityErr400_MalformedJSONErr400_InvalidRequestBody"
	_ErrorCode_name_2 = "Err401_InvalidCredentialsErr401_AuthorizationHeaderMissingErr401_AuthorizationHeaderInvalidErr401_UserNotFound"
	_ErrorCode_name_3 = "Err403_CannotToDeleteErr403_CannotEditPhone"
	_ErrorCode_name_4 = "Err404_PlayerStatsNotFoundErr404_UserOrPhoneNotFoundErr404_AccountNotFound"
	_ErrorCode_name_5 = "Err429_EditRequestTimedOut"
	_ErrorCode_name_6 = "Err500_UnableToDeleteErr500_UnableToEditPhoneErr500_UnableToRegisterErr500_UnknownError"
	_ErrorCode_name_7 = "Err503_DataBaseOnDeleteErr503_DataBaseOnPhoneEdit"
)

var (
	_ErrorCode_index_1 = [...]uint8{0, 25, 50, 78, 103, 132, 158, 194, 214, 239}
	_ErrorCode_index_2 = [...]uint8{0, 25, 58, 91, 110}
	_ErrorCode_index_3 = [...]uint8{0, 21, 43}
	_ErrorCode_index_4 = [...]uint8{0, 26, 52, 74}
	_ErrorCode_index_6 = [...]uint8{0, 21, 45, 68, 87}
	_ErrorCode_index_7 = [...]uint8{0, 23, 49}
)

func (i ErrorCode) String() string {
	switch {
	case i == 2070001:
		return _ErrorCode_name_0
	case 4000001 <= i && i <= 4000009:
		i -= 4000001
		return _ErrorCode_name_1[_ErrorCode_index_1[i]:_ErrorCode_index_1[i+1]]
	case 4010001 <= i && i <= 4010004:
		i -= 4010001
		return _ErrorCode_name_2[_ErrorCode_index_2[i]:_ErrorCode_index_2[i+1]]
	case 4030001 <= i && i <= 4030002:
		i -= 4030001
		return _ErrorCode_name_3[_ErrorCode_index_3[i]:_ErrorCode_index_3[i+1]]
	case 4040001 <= i && i <= 4040003:
		i -= 4040001
		return _ErrorCode_name_4[_ErrorCode_index_4[i]:_ErrorCode_index_4[i+1]]
	case i == 4290001:
		return _ErrorCode_name_5
	case 5000001 <= i && i <= 5000004:
		i -= 5000001
		return _ErrorCode_name_6[_ErrorCode_index_6[i]:_ErrorCode_index_6[i+1]]
	case 5030001 <= i && i <= 5030002:
		i -= 5030001
		return _ErrorCode_name_7[_ErrorCode_index_7[i]:_ErrorCode_index_7[i+1]]
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
