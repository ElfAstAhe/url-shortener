package repository

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortUriRepository interface {
	GetById(id string) (*_model.ShortUri, error)
	GetByKey(key string) (*_model.ShortUri, error)
	Create(shortUri *_model.ShortUri) (*_model.ShortUri, error)
}

func NewShortUriRepository(dbConfig *_cfg.DbConfig) ShortUriRepository {
	// check in future (for next dev iteration)
	if dbConfig == nil || dbConfig.Kind == _cfg.DbKindInMemory {
		return NewShortUriInMemRepo()
	} else if dbConfig.Kind == _cfg.DbKindPostgres {
		return NewShortUriInMemRepo()
	}

	return NewShortUriInMemRepo()
}
