package repository

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIRepository interface {
	Get(id string) (*_model.ShortURI, error)
	GetByKey(key string) (*_model.ShortURI, error)
	Create(shortURI *_model.ShortURI) (*_model.ShortURI, error)
	BatchCreate(batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error)
	ListAllByUser(userID string) ([]*_model.ShortURI, error)
}

func NewShortURIRepository(db _db.DB) (ShortURIRepository, error) {
	if db != nil && db.GetDBKind() == _cfg.DBKindPostgres {
		return newShortURIPgRepo(db)
	}

	return newShortURIInMemRepo(db)
}
