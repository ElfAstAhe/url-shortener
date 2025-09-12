package db

import (
	"database/sql"
	"io"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
)

type DB interface {
	GetDB() *sql.DB
	GetDBKind() string
	GetDsn() string
}

func NewDB(kind string, dsn string) (DB, error) {
	if kind == _cfg.DBKindPostgres {
		return newPostgresqlDB(kind, dsn)
	}

	return newInMemoryDB(_cfg.DBKindInMemory)
}

func CloseDB(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
