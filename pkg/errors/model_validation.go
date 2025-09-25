package errors

import (
	"fmt"
)

type AppModelValidationError struct {
	ID     string
	Entity string
	Err    error
	Msg    string
}

func NewAppModelValidationError(entity string, err error) *AppModelValidationError {
	return NewAppModelValidationFullError("", entity, err, "")
}

func NewAppModelValidationFullError(id, entity string, err error, msg string) *AppModelValidationError {
	return &AppModelValidationError{
		ID:     id,
		Entity: entity,
		Err:    err,
		Msg:    msg,
	}
}

func (e *AppModelValidationError) Error() string {
	switch {
	// full
	case e.Msg != "" && e.Err != nil && e.ID != "" && e.Entity != "":
		return fmt.Sprintf("model [%s] with id [%s] validation error with err [%s] ext messasge [%s]", e.Entity, e.ID, e.Err.Error(), e.Msg)
	// partial
	case e.Msg != "" && e.Err == nil && e.ID != "" && e.Entity != "":
		return fmt.Sprintf("model [%s] with id [%s] validation error with ext message [%s]", e.Entity, e.ID, e.Msg)
	case e.Msg == "" && e.Err != nil && e.ID != "" && e.Entity != "":
		return fmt.Sprintf("model [%s] with id [%s] validation error with err [%s]", e.Entity, e.ID, e.Err.Error())
	case e.Entity != "" && e.Err != nil:
		return fmt.Sprintf("model [%s] validation error with err [%s]", e.ID, e.Err.Error())
	}

	return "model validation error"
}

func (e *AppModelValidationError) Unwrap() error {
	return e.Err
}
