package db

import (
	"database/sql"
	"io"

	"github.com/ElfAstAhe/url-shortener/internal/config"
)

type DB interface {
	GetDB() *sql.DB
	GetDBKind() string
	GetDsn() string
}

func NewDB(kind string, dsn string) (DB, error) {
	if kind == config.DBKindPostgres {
		return newPostgresqlDB(kind, dsn)
	}

	return newInMemoryDB(config.DBKindInMemory)
}

func CloseDB(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
