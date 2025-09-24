package repository

import (
	"context"
	"database/sql"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIUserRepository interface {
	Get(ctx context.Context, ID string) (*_model.ShortURIUser, error)
	GetByUnique(ctx context.Context, userID string, shortURIID string) (*_model.ShortURIUser, error)
	ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURIUser, error)
	ListAllByShortURI(ctx context.Context, shortURIID string) ([]*_model.ShortURIUser, error)
	Create(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	CreateTran(ctx context.Context, tx *sql.Tx, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	Change(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	//Delete(ctx context.Context, ID string) error
	//BatchDelete(ctx context.Context, ids []string) error
}

func NewShortURIUserRepository(db _db.DB) (ShortURIUserRepository, error) {
	if db != nil && db.GetDBKind() == _cfg.DBKindPostgres {
		return newShortURIUserPgRepo(db)
	}

	return newShortURIUserInMemRepo(db)
}
