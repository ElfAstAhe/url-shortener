package repository

import (
	"context"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/config/db"
	"github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIRepository interface {
	Get(ctx context.Context, id string) (*model.ShortURI, error)
	GetByKey(ctx context.Context, key string) (*model.ShortURI, error)
	GetByKeyUser(ctx context.Context, userID string, key string) (*model.ShortURI, error)
	ListAllByUser(ctx context.Context, userID string) ([]*model.ShortURI, error)
	ListAllByKeys(ctx context.Context, keys []string) ([]*model.ShortURI, error)
	Create(ctx context.Context, userID string, entity *model.ShortURI) (*model.ShortURI, error)
	BatchCreate(ctx context.Context, userID string, batch map[string]*model.ShortURI) (map[string]*model.ShortURI, error)
	Delete(ctx context.Context, ID string, userID string) error
	BatchDeleteByKeys(ctx context.Context, userID string, keys []string) error
}

func NewShortURIRepository(db db.DB) (ShortURIRepository, error) {
	if db != nil && db.GetDBKind() == config.DBKindPostgres {
		return newShortURIPgRepo(db)
	}

	return newShortURIInMemRepo(db)
}
