package repository

import (
	"database/sql"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIUserRepository interface {
	Get(ID string) (*_model.ShortURIUser, error)
	GetByUnique(userID string, shortURIID string) (*_model.ShortURIUser, error)
	ListAllByUser(userID string) ([]*_model.ShortURIUser, error)
	ListAllByShortURI(shortURIID string) ([]*_model.ShortURIUser, error)
	Create(entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	CreateTran(tx *sql.Tx, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
}

func NewShortURIUserRepository(db _db.DB) (ShortURIUserRepository, error) {
	if db != nil && db.GetDBKind() == _cfg.DBKindPostgres {
		return newShortURIUserPgRepo(db)
	}

	return newShortURIUserInMemRepo(db)
}
