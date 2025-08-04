package errors

import "fmt"

type NotFoundError struct {
	key string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("Short URI with key [%s] not found", n.key)
}

func (n NotFoundError) RuntimeError() {
}

func NewNotFoundError(key string) *NotFoundError {
	return &NotFoundError{key: key}
}
