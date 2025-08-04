package errors

import "fmt"

type InvalidOriginalUrlError struct {
	OriginalUrl string
}

func (i InvalidOriginalUrlError) Error() string {
	return fmt.Sprintf("OriginalUrl ['%s'] is invalid", i.OriginalUrl)
}

func (i InvalidOriginalUrlError) RuntimeError() {
}

func NewInvalidOriginalUrlError(originalUrl string) *InvalidOriginalUrlError {
	return &InvalidOriginalUrlError{OriginalUrl: originalUrl}
}
