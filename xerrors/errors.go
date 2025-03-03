package xerrors

import "fmt"

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func New(code int, msg string) error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}
