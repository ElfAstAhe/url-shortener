package db

import (
	"database/sql"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type inMemoryDB struct {
	ShortURI map[string]*_model.ShortURI
	DBKind   string
	Config   *_cfg.DBConfig
}

func newInMemoryDB(kind string, config *_cfg.DBConfig) (*inMemoryDB, error) {
	return &inMemoryDB{
		ShortURI: make(map[string]*_model.ShortURI),
		DBKind:   kind,
		Config:   config,
	}, nil
}

// Closer

func (db *inMemoryDB) Close() error {
	clear(db.ShortURI)

	return nil
}

// ========

// DB

func (db *inMemoryDB) GetDB() *sql.DB {
	return nil
}

func (db *inMemoryDB) GetDBKind() string {
	return db.DBKind
}

func (db *inMemoryDB) GetConfig() *_cfg.DBConfig {
	return db.Config
}

// ========

// InMemoryCache

func (db *inMemoryDB) GetShortURICache() map[string]*_model.ShortURI {
	return db.ShortURI
}

// ========
