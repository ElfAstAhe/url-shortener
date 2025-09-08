package db

import (
	"database/sql"
	"io"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
)

type DB interface {
	GetDB() *sql.DB
	GetDBKind() string
	GetConfig() *_cfg.DBConfig
}

func NewDB(kind string, config *_cfg.DBConfig) (DB, error) {
	if kind == _cfg.DBKindPostgres {
		return newPostgresqlDB(kind, config)
	}

	return newInMemoryDB(_cfg.DBKindInMemory, config)
}

func CloseDB(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
