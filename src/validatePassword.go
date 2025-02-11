package src

import "regexp"

func ValidatePassword(password string) bool {
	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`\d`)
	return lower.MatchString(password) && upper.MatchString(password) && digit.MatchString(password) && len(password) >= 8
}
