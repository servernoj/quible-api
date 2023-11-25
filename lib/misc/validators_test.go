package misc

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestPhoneValidator_InvalidInput(t *testing.T) {
	validate := validator.New()
	RegisterValidators(validate)

	testCases := []struct {
		input    interface{}
		expected bool
	}{
		{"1234567890", true},          // Valid phone number
		{"123-456-7890", true},        // Valid phone number with dashes
		{"123 456 7890", true},        // Valid phone number with spaces
		{"123456789", false},          // Invalid phone number (less than 10 digits)
		{"12345678901", true},         // Valid phone number (more than 10 digits)
		{1234567890, false},           // Invalid input type (not a string)
		{[]byte("1234567890"), false}, // Invalid input type (not a string)
	}

	for _, tc := range testCases {
		t.Run("Input_"+reflect.TypeOf(tc.input).String(), func(t *testing.T) {
			err := validate.Var(tc.input, "phone")
			if (err == nil) != tc.expected {
				t.Errorf("Expected phoneValidator to return %v for input %v, but got error: %v", tc.expected, tc.input, err)
			}
		})
	}
}
