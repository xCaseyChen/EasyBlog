package utils

import (
	"regexp"
)

var (
	upperRegex  = regexp.MustCompile("[A-Z]")
	lowerRegex  = regexp.MustCompile("[a-z]")
	numberRegex = regexp.MustCompile("[0-9]")
)

func ValidatePasswordCompliance(s string) bool {
	return len(s) >= 8 && upperRegex.MatchString(s) && lowerRegex.MatchString(s) && numberRegex.MatchString(s)
}
