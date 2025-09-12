package db

import (
	"database/sql"

	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type inMemoryDB struct {
	ShortURI map[string]*_model.ShortURI
	DBKind   string
}

var inMemDB *inMemoryDB

func newInMemoryDB(kind string) (*inMemoryDB, error) {
	if inMemDB != nil {
		return inMemDB, nil
	}

	inMemDB = &inMemoryDB{
		ShortURI: make(map[string]*_model.ShortURI),
		DBKind:   kind,
	}

	return inMemDB, nil
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

func (db *inMemoryDB) GetDsn() string {
	return ""
}

// ========

// InMemoryCache

func (db *inMemoryDB) GetShortURICache() map[string]*_model.ShortURI {
	return db.ShortURI
}

// ========
