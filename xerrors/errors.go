package xerrors

import "fmt"

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// Detail is a short hand function to create a New error by adding err.Error() to e.Msg
func (e *Error) Detail(err error) error {
	return New(
		e.Code,
		fmt.Sprintf("%s: %s", e.Msg, err.Error()),
	)
}

func New(code int, msg string) error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}
