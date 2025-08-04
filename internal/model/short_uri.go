package model

import (
	"github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
	"net/url"
)

type ShortUri struct {
	//
	Id          string  `db:"id"`
	OriginalUrl url.URL `db:"original_url"`
	Key         string  `db:"key"`
	TechData
}

func NewShortUri(originalUrl string, key string) (*ShortUri, error) {
	origUrl, err := url.Parse(originalUrl)
	if err != nil {
		return nil, errors.NewInvalidOriginalUrlError(originalUrl)
	}
	uri := &ShortUri{
		OriginalUrl: *origUrl,
		Key:         key,
		Id:          uuid.New().String(),
	}
	return uri, nil
}
