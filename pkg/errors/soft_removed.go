package errors

import (
	"fmt"
)

type AppSoftRemovedError struct {
	Entity string
	Err    error
}

var AppSoftRemoved *AppSoftRemovedError

func NewAppSoftRemovedError(entity string, err error) *AppSoftRemovedError {
	return &AppSoftRemovedError{
		Entity: entity,
		Err:    err,
	}
}

func (e *AppSoftRemovedError) Error() string {
	return fmt.Sprintf("entity [%s] soft removed", e.Entity)
}

func (e *AppSoftRemovedError) Unwrap() error { return e.Err }
