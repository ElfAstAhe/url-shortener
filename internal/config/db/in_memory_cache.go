package db

import (
	"sync"

	"github.com/ElfAstAhe/url-shortener/internal/model"
)

type InMemoryCache interface {
	GetShortURIRWMutex() *sync.RWMutex
	GetShortURIUserRWMutex() *sync.RWMutex
	GetShortURICache() map[string]*model.ShortURI
	GetShortURIUserCache() map[string]*model.ShortURIUser
	GetShortURIAuditCache() map[string]*model.ShortURIAudit
}
