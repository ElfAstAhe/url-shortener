package errors

import "fmt"

type AppNotFoundError struct {
	key string
}

var AppNotFound *AppNotFoundError

func NewAppNotFoundError(key string) *AppNotFoundError {
	return &AppNotFoundError{key: key}
}

func (n *AppNotFoundError) Error() string {
	return fmt.Sprintf("Short URI with key [%s] not found", n.key)
}
