// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package controller

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Err400_UnknownError-4003001]
	_ = x[Err400_InvalidRequestBody-4003002]
	_ = x[Err401_UnknownError-4013001]
	_ = x[Err404_UnknownError-4043001]
	_ = x[Err424_UnknownError-4243001]
	_ = x[Err424_PostmarkSendEmail-4243002]
	_ = x[Err500_UnknownError-5003001]
}

const (
	_ErrorCode_name_0 = "Err400_UnknownErrorErr400_InvalidRequestBody"
	_ErrorCode_name_1 = "Err401_UnknownError"
	_ErrorCode_name_2 = "Err404_UnknownError"
	_ErrorCode_name_3 = "Err424_UnknownErrorErr424_PostmarkSendEmail"
	_ErrorCode_name_4 = "Err500_UnknownError"
)

var (
	_ErrorCode_index_0 = [...]uint8{0, 19, 44}
	_ErrorCode_index_3 = [...]uint8{0, 19, 43}
)

func (i ErrorCode) String() string {
	switch {
	case 4003001 <= i && i <= 4003002:
		i -= 4003001
		return _ErrorCode_name_0[_ErrorCode_index_0[i]:_ErrorCode_index_0[i+1]]
	case i == 4013001:
		return _ErrorCode_name_1
	case i == 4043001:
		return _ErrorCode_name_2
	case 4243001 <= i && i <= 4243002:
		i -= 4243001
		return _ErrorCode_name_3[_ErrorCode_index_3[i]:_ErrorCode_index_3[i+1]]
	case i == 5003001:
		return _ErrorCode_name_4
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
