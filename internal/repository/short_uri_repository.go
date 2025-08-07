package repository

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIRepository interface {
	GetByID(id string) (*_model.ShortURI, error)
	GetByKey(key string) (*_model.ShortURI, error)
	Create(shortURI *_model.ShortURI) (*_model.ShortURI, error)
}

func NewShortURIRepository(dbConfig *_cfg.DBConfig) ShortURIRepository {
	// check in future (for next dev iteration)
	if dbConfig == nil || dbConfig.Kind == _cfg.DBKindInMemory {
		return NewShortURIInMemRepo()
	} else if dbConfig.Kind == _cfg.DBKindPostgres {
		return NewShortURIInMemRepo()
	}

	return NewShortURIInMemRepo()
}
