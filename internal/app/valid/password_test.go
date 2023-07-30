package valid

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type isValidPasswordTest struct {
	description    string
	password       string
	expectedResult bool
}

func TestIsValidPassword(t *testing.T) {
	tests := []isValidPasswordTest{
		{
			description:    "valid password",
			password:       "qwerty_123",
			expectedResult: true,
		},
		{
			description:    "too short password",
			password:       "zzz_0",
			expectedResult: false,
		},
		{
			description:    "too long password",
			password:       "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz_123",
			expectedResult: false,
		},
		{
			description:    "password with no latin letters",
			password:       "123_123",
			expectedResult: false,
		},
		{
			description:    "password with no digits",
			password:       "abcdefg_ABCDEFG",
			expectedResult: false,
		},
		{
			description:    "password with no special symbols",
			password:       "abcdefg12345",
			expectedResult: false,
		},
		{
			description:    "password with not allowed symbols",
			password:       "абвгд_12345",
			expectedResult: false,
		},
	}
	for _, test := range tests {
		assert.Equal(t, IsValidPassword(test.password), test.expectedResult)
	}
}
