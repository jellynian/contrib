package validate

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

func (f Form) Errors() []*Error {
	return f.errors
}

func (f *Form) IsExactly() bool {
	return len(f.errors) == 0
}
