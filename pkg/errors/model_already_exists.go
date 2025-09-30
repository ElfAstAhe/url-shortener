package errors

import "fmt"

type AppModelAlreadyExistsError struct {
	ID     string
	Entity string
}

var AppModelAlreadyExists *AppModelAlreadyExistsError

func NewAppModelAlreadyExistsError(ID string, entity string) *AppModelAlreadyExistsError {
	return &AppModelAlreadyExistsError{
		ID:     ID,
		Entity: entity,
	}
}

func (e *AppModelAlreadyExistsError) Error() string {
	return fmt.Sprintf("model [%s] already exists with id [%s]", e.Entity, e.ID)
}
