package errors

import "fmt"

type AppInvalidArgumentError struct {
	argument string
}

var AppInvalidArgument *AppInvalidArgumentError

func NewAppInvalidArgument(argument string) *AppInvalidArgumentError {
	return &AppInvalidArgumentError{
		argument: argument,
	}
}

func (e *AppInvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument error: [%s]", e.argument)
}
