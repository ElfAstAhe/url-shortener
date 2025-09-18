package db

import (
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type InMemoryCache interface {
	GetShortURICache() map[string]*_model.ShortURI
	GetShortURIUserCache() map[string]*_model.ShortURIUser
	GetShortURIAuditCache() map[string]*_model.ShortURIAudit
}
