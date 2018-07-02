package validate

import (
	"regexp"
	"strings"
)

var (
	regexEmail    = "[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?"
	regexpIsEmail = regexp.MustCompile(regexEmail)
)

type Error struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

type Form struct {
	errors []*Error
}

func (f *Form) AppendError(err *Error) {
	f.errors = append(f.errors, err)
}

func (f *Form) Errors() []*Error {
	return f.errors
}

func (f *Form) IsExact() bool {
	return len(f.errors) == 0
}

func IsEmail(val string) bool {
	return regexpIsEmail.MatchString(val)
}

func IsEmpty(val string) bool {
	return len(strings.TrimSpace(val)) == 0
}
