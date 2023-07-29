package valid

import "strings"

const allowedNicknameSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func isAllowed(c rune, allowedSymbols string) bool {
	for _, s := range allowedSymbols {
		if c == s {
			return true
		}
	}
	return false
}

var obscenities = map[string]struct{}{
	"fuck":   {},
	"nigger": {},
}

// IsValidNickname checks if nickname has valid len, contains only allowed
// symbols and doesn't contain obscenities
func IsValidNickname(nickname string) bool {

	// check if nickname have valid len
	if !(len(nickname) >= 4 && len(nickname) <= 25) {
		return false
	}

	// check if nickname contains only allowed symbols
	for _, s := range nickname {
		if !isAllowed(s, allowedNicknameSymbols) {
			return false
		}
	}

	// check if nickname doesn't contain obscenities
	lowerNickname := strings.ToLower(nickname)
	for word := range obscenities {
		if strings.Contains(lowerNickname, word) {
			return false
		}
	}

	return true
}
