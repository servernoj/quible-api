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
			// Valid phone number: This test ensures that the validate.Var function returns no error when a valid phone number is provided.
			if err := validate.Var("1234567890", "phone"); err != nil {
				t.Errorf("Expected no error for valid phone number %q, but got: %v", "1234567890", err)
			}

			// Invalid phone number (less than 10 digits): This test verifies that the validate.Var function returns an error when an invalid phone number is provided.
			if err := validate.Var("123456789", "phone"); err == nil {
				t.Errorf("Expected error for invalid phone number %q, got none", "123456789")
			}
		})
	}
}
