package db

import (
	"sync"

	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type InMemoryCache interface {
	GetShortURIRWMutex() *sync.RWMutex
	GetShortURIUserRWMutex() *sync.RWMutex
	GetShortURICache() map[string]*_model.ShortURI
	GetShortURIUserCache() map[string]*_model.ShortURIUser
	GetShortURIAuditCache() map[string]*_model.ShortURIAudit
}
