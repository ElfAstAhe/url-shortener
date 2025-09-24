package errors

import "fmt"

type AppAuthCookieAbsentError struct {
	msg string
	Err error
}

var AppAuthCookieAbsent *AppAuthCookieAbsentError

func NewAppAuthCookieAbsentError(msg string, err error) *AppAuthCookieAbsentError {
	return &AppAuthCookieAbsentError{
		msg: msg,
		Err: err,
	}
}

func (a *AppAuthCookieAbsentError) Error() string {
	switch {
	case a.msg != "" && a.Err != nil:
		return fmt.Sprintf("auth cookie absent with msg [%s] and err [%v]", a.msg, a.Err)
	case a.msg == "":
		return fmt.Sprintf("auth cookie absent with msg [%s]", a.msg)
	case a.Err != nil:
		return fmt.Sprintf("auth cookie absent with err [%v]", a.Err)
	}

	return "auth cookie absent"
}

func (a *AppAuthCookieAbsentError) Unwrap() error {
	return a.Err
}
