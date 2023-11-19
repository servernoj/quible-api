package misc

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func PickFields(data any, fields ...string) map[string]any {
	bytes, _ := json.Marshal(&data)
	fullMap := make(map[string]any)
	json.Unmarshal(bytes, &fullMap)
	result := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := fullMap[f]; ok {
			result[f] = v
		}
	}
	return result
}

type ErrorFields struct {
	CheckSome    func(...string) bool
	CheckAll     func(...string) bool
	GetAllFields func() []string
}

// Function parses validation error and returns a CLOSURE with 3 functions:
//
//	getAllFields() -- reports all problem fields as a slice of strings formatted as `field:problem`
//	checkSome() -- checks if SOME of the listed fields have been reported to have valiadtion errors
//	checkAll() -- checks if ALL of the listed fields have been reported to have valiadtion errors
func ParseValidationError(err error) ErrorFields {
	set := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			key := fe.StructField()
			set[key] = fe.Tag()
		}
	}
	idx := 0
	allFields := make([]string, len(set))
	for key, value := range set {
		allFields[idx] = key + ":" + value
		idx++
	}

	getAllFields := func() []string {
		return allFields
	}

	checkSome := func(keys ...string) bool {
		flag := false
		for _, key := range keys {
			_, ok := set[key]
			if ok {
				flag = true
				break
			}
		}
		return flag
	}

	checkAll := func(keys ...string) bool {
		flag := false
		for _, key := range keys {
			_, ok := set[key]
			if !ok {
				flag = true
				break
			}
		}
		return !flag
	}
	return ErrorFields{
		CheckSome:    checkSome,
		CheckAll:     checkAll,
		GetAllFields: getAllFields,
	}
}
