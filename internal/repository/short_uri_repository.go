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

func NewShortURIRepository(dbKind string, dbConfig *_cfg.DBConfig) ShortURIRepository {
	// check in future (for next dev iteration)
	if dbKind == "" || dbKind == _cfg.DefaultDBKind {
		return NewShortURIInMemRepo(dbKind, dbConfig)
	} else if dbKind == _cfg.DBKindPostgres {
		return NewShortURIInMemRepo(dbKind, dbConfig)
	}

	return NewShortURIInMemRepo(dbKind, dbConfig)
}
