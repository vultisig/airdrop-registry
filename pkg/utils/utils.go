package utils

import "regexp"

func IsValidHex(s string) bool {
	re := regexp.MustCompile("^[0-9a-fA-F]{64}$")
	return re.MatchString(s)
}
