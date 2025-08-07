package model

import (
	"net/url"

	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
)

type ShortURI struct {
	//
	ID          string  `db:"id"`
	OriginalURL url.URL `db:"original_url"`
	Key         string  `db:"key"`
	TechData
}

func NewShortURI(originalURL string, key string) (*ShortURI, error) {
	origURL, err := url.Parse(originalURL)
	if err != nil {
		return nil, _err.NewInvalidOriginalURLError(originalURL)
	}
	uri := &ShortURI{
		OriginalURL: *origURL,
		Key:         key,
		ID:          uuid.New().String(),
	}
	return uri, nil
}
