package errors

import "fmt"

type AppInvalidOriginalURLError struct {
	OriginalURL string
}

var AppInvalidOriginalURL *AppInvalidOriginalURLError

func NewInvalidOriginalURLError(originalURL string) *AppInvalidOriginalURLError {
	return &AppInvalidOriginalURLError{OriginalURL: originalURL}
}

func (i *AppInvalidOriginalURLError) Error() string {
	return fmt.Sprintf("OriginalURL ['%s'] is invalid", i.OriginalURL)
}
