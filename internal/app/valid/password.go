package valid

import "unicode"

const allowedPasswordSpecial = "!@#$%^&*()_+-=,.<>?;:{}[]"

// IsValidPassword checks if password has valid len and contains only allowed
// symbols and has all of letter, digit and special symbol
func IsValidPassword(password string) bool {

	// check if password has valid len
	if !(len(password) >= 6 && len(password) <= 50) {
		return false
	}

	// check if password contains only allowed symbols and has all of letter, digit and special symbol
	var hasLetter, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.Is(unicode.Latin, c):
			hasLetter = true
		case unicode.IsDigit(c):
			hasDigit = true
		case isAllowed(c, allowedPasswordSpecial):
			hasSpecial = true
		default:
			return false
		}
	}
	if !(hasLetter && hasDigit && hasSpecial) {
		return false
	}

	return true
}
