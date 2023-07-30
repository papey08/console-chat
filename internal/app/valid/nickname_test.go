package valid

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type isValidNicknameTest struct {
	description    string
	nickname       string
	expectedResult bool
}

func TestIsValidNickname(t *testing.T) {
	tests := []isValidNicknameTest{
		{
			description:    "valid nickname",
			nickname:       "papey08",
			expectedResult: true,
		},
		{
			description:    "too short nickname",
			nickname:       "zzz",
			expectedResult: false,
		},
		{
			description:    "too long nickname",
			nickname:       "abcdefghijklmnopqrstuvwxyz",
			expectedResult: false,
		},
		{
			description:    "nickname with wrong symbols",
			nickname:       "papey08_!@#$",
			expectedResult: false,
		},
		{
			description:    "obscenity nickname",
			nickname:       "fuck",
			expectedResult: false,
		},
		{
			description:    "nickname with obscenity",
			nickname:       "fuck_you",
			expectedResult: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, IsValidNickname(test.nickname), test.expectedResult)
	}
}
