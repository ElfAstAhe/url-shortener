package model

import (
	"errors"
	"net/url"

	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
)

type ShortURI struct {
	ID          string     `db:"id" json:"id"`
	OriginalURL *CustomURL `db:"original_url" json:"original_url"`
	Key         string     `db:"key" json:"key"`
}

func NewShortURI(originalURL string, key string) (*ShortURI, error) {
	return NewShortURIFull(uuid.New().String(), originalURL, key)
}

func NewShortURIFull(ID string, originalURL string, key string) (*ShortURI, error) {
	origURL, err := url.Parse(originalURL)
	if err != nil {
		return nil, _err.NewInvalidOriginalURLError(originalURL)
	}

	return &ShortURI{
		ID:          ID,
		OriginalURL: &CustomURL{origURL},
		Key:         key,
	}, nil
}

func ValidateShortURI(entity *ShortURI) error {
	if entity == nil {
		return errors.New("entity is nil")
	}
	if entity.Key == "" {
		return errors.New("key is required")
	}
	if entity.OriginalURL == nil || entity.OriginalURL.URL == nil || entity.OriginalURL.URL.String() == "" {
		return errors.New("original_url is required")
	}

	return nil
}
