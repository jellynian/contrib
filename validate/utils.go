package validate

import (
	"strings"
)

func IsEmail(val string) bool {
	return regexpIsEmail.MatchString(val)
}

func IsEmpty(val string) bool {
	return len(strings.Trim(val, " ")) == 0
}
