package repository

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
)

type DBConnCheckRepository interface {
	CheckDBConn() error
}

func NewDBConnCheckRepository(db _db.DB) (DBConnCheckRepository, error) {
	if db.GetDBKind() == _cfg.DBKindPostgres {
		return newDBConnCheckPgRepo(db)
	}

	return newDBConnCheckImMemRepo()
}
