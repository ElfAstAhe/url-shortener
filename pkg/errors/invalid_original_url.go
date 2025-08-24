package errors

import "fmt"

type InvalidOriginalURLError struct {
	OriginalURL string
}

func (i InvalidOriginalURLError) Error() string {
	return fmt.Sprintf("OriginalURL ['%s'] is invalid", i.OriginalURL)
}

func (i InvalidOriginalURLError) RuntimeError() {
}

func NewInvalidOriginalURLError(originalURL string) *InvalidOriginalURLError {
	return &InvalidOriginalURLError{OriginalURL: originalURL}
}
