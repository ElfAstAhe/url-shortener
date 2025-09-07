package db

import (
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type InMemoryCache interface {
	GetShortURICache() map[string]*_model.ShortURI
}
