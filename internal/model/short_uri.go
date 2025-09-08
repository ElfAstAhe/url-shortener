package model

import (
	"net/url"
	"time"

	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
)

type ShortURI struct {
	//
	ID          string  `db:"id" json:"id"`
	OriginalURL url.URL `db:"original_url" json:"original_url"`
	Key         string  `db:"key" json:"key"`
	TechData    `json:"tech_data,omitempty"`
}

func NewShortURI(originalURL string, key string) (*ShortURI, error) {
	return NewShortURIFull(uuid.New().String(), originalURL, key, &TechData{
		CreateUser: "unknown",
		Created:    time.Now(),
		UpdateUser: "unknown",
		Updated:    time.Now(),
	})
}

func NewShortURIFull(ID string, originalURL string, key string, techData *TechData) (*ShortURI, error) {
	origURL, err := url.Parse(originalURL)
	if err != nil {
		return nil, _err.NewInvalidOriginalURLError(originalURL)
	}

	return NewShortURIComplete(ID, origURL, key, techData), nil
}

func NewShortURIComplete(ID string,
	origURL *url.URL,
	key string,
	techData *TechData) *ShortURI {
	return &ShortURI{
		ID:          ID,
		OriginalURL: *origURL,
		Key:         key,
		TechData:    *techData,
	}
}
