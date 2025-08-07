package model

import (
	"net/url"

	"github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
)

type ShortURI struct {
	//
	Id          string  `db:"id"`
	OriginalURL url.URL `db:"original_url"`
	Key         string  `db:"key"`
	TechData
}

func NewShortUri(originalURL string, key string) (*ShortURI, error) {
	origUrl, err := url.Parse(originalURL)
	if err != nil {
		return nil, errors.NewInvalidOriginalURLError(originalURL)
	}
	uri := &ShortURI{
		OriginalURL: *origUrl,
		Key:         key,
		Id:          uuid.New().String(),
	}
	return uri, nil
}
