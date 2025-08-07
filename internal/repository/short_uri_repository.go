package repository

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIRepository interface {
	GetById(id string) (*_model.ShortURI, error)
	GetByKey(key string) (*_model.ShortURI, error)
	Create(shortUri *_model.ShortURI) (*_model.ShortURI, error)
}

func NewShortUriRepository(dbConfig *_cfg.DBConfig) ShortURIRepository {
	// check in future (for next dev iteration)
	if dbConfig == nil || dbConfig.Kind == _cfg.DBKindInMemory {
		return NewShortUriInMemRepo()
	} else if dbConfig.Kind == _cfg.DBKindPostgres {
		return NewShortUriInMemRepo()
	}

	return NewShortUriInMemRepo()
}
