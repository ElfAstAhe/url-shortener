package repository

import (
	"context"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIRepository interface {
	Get(ctx context.Context, id string) (*_model.ShortURI, error)
	GetByKey(ctx context.Context, key string) (*_model.ShortURI, error)
	GetByKeyUser(ctx context.Context, userID string, key string) (*_model.ShortURI, error)
	Create(ctx context.Context, entity *_model.ShortURI) (*_model.ShortURI, error)
	BatchCreate(ctx context.Context, batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error)
	ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURI, error)
}

func NewShortURIRepository(db _db.DB) (ShortURIRepository, error) {
	if db != nil && db.GetDBKind() == _cfg.DBKindPostgres {
		return newShortURIPgRepo(db)
	}

	return newShortURIInMemRepo(db)
}
