package repository

import (
	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/config/db"
)

type DBConnCheckRepository interface {
	CheckDBConn() error
}

func NewDBConnCheckRepository(appDb db.DB) (DBConnCheckRepository, error) {
	if appDb.GetDBKind() == config.DBKindPostgres {
		return newDBConnCheckPgRepo(appDb)
	}

	return newDBConnCheckImMemRepo()
}
