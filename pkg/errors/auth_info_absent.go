package errors

import "fmt"

type AppAuthInfoAbsentError struct {
	msg string
	Err error
}

func (a AppAuthInfoAbsentError) Error() string {
	switch {
	case a.msg != "" && a.Err != nil:
		return fmt.Sprintf("auth user info absent with msg [%s] and err [%v]", a.msg, a.Err)
	case a.msg == "":
		return fmt.Sprintf("auth user info absent with msg [%s]", a.msg)
	case a.Err != nil:
		return fmt.Sprintf("auth user info absent with err [%v]", a.Err)
	}

	return fmt.Sprint("auth user info absent")
}

func (a AppAuthInfoAbsentError) Unwrap() error {
	return a.Err
}

func NewAppAuthInfoAbsentError(msg string, err error) *AppAuthInfoAbsentError {
	return &AppAuthInfoAbsentError{
		msg: msg,
		Err: err,
	}
}
