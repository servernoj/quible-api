package misc

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestPhoneValidator_InvalidInput(t *testing.T) {
	validate := validator.New()
	RegisterValidators(validate)

	testCases := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "Valid phone number: 1234567890",
			input:    "1234567890",
			expected: true,
		},
		{
			name:     "Valid phone number with dashes: 123-456-7890",
			input:    "123-456-7890",
			expected: true,
		},
		{
			name:     "Valid phone number with spaces: 123 456 7890",
			input:    "123 456 7890",
			expected: true,
		},
		{
			name:     "Invalid phone number (less than 10 digits): 123456789",
			input:    "123456789",
			expected: false,
		},
		{
			name:     "Valid phone number (more than 10 digits): 12345678901",
			input:    "12345678901",
			expected: true,
		},
		{
			name:     "Invalid input type (not a string): integer",
			input:    1234567890,
			expected: false,
		},
		{
			name:     "Invalid input type (not a string): byte array",
			input:    []byte("1234567890"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Var(tc.input, "phone")
			if got, wanted := err == nil, tc.expected; got != wanted {
				t.Errorf("Got %t, Expected %t for test case: %s", got, wanted, tc.name)
			}
		})
	}
}
