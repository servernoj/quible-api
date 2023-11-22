package misc

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func PickFields(data any, fields ...string) map[string]any {
	bytes, _ := json.Marshal(&data)
	fullMap := make(map[string]any)
	_ = json.Unmarshal(bytes, &fullMap)
	result := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := fullMap[f]; ok {
			result[f] = v
		}
	}
	return result
}

type ErrorFields struct {
	IsValidationError bool
	CheckSome         func(...string) bool
	CheckAll          func(...string) bool
	GetAllFields      func() []string
}

// Function parses validation error and returns a CLOSURE with 3 functions and one flag:
//
//	getAllFields() -- reports all problem fields as a slice of strings formatted as `field:problem`
//	checkSome() -- checks if SOME of the listed fields have been reported to have validation errors
//	checkAll() -- checks if ALL of the listed fields have been reported to have validation errors
func ParseValidationError(err error) ErrorFields {
	set := make(map[string]string)
	IsValidationError := false

	// Added error handling to check the type of input error
	if _, ok := err.(validator.ValidationErrors); !ok {
		return ErrorFields{IsValidationError: false} // Return an appropriate error message if invalid input
	}

	IsValidationError = true
	validationErrors := err.(validator.ValidationErrors)
	for _, fe := range validationErrors {
		key := fe.StructField()
		set[key] = fe.Tag()
	}

	// Prepopulate 'allFields' slice with formatted strings instead of constructing each element separately
	allFields := make([]string, len(set))
	idx := 0
	for key, value := range set {
		allFields[idx] = key + ":" + value
		idx++
	}

	getAllFields := func() []string {
		return allFields
	}

	checkSome := func(keys ...string) bool {
		for _, key := range keys {
			if _, ok := set[key]; ok {
				return true
			}
		}
		return false
	}

	// Removed redundant flag in 'checkAll' function
	checkAll := func(keys ...string) bool {
		for _, key := range keys {
			if _, ok := set[key]; !ok {
				return false
			}
		}
		return true // Simplified code and avoided an additional 'flag' variable
	}

	return ErrorFields{
		IsValidationError: IsValidationError,
		CheckSome:         checkSome,
		CheckAll:          checkAll,
		GetAllFields:      getAllFields,
	}
}
