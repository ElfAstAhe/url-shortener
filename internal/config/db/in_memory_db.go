package db

import (
	"database/sql"
	"sync"

	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type inMemoryDB struct {
	DBKind string

	// ShortURI map key is short uri entity id attribute
	ShortURI map[string]*_model.ShortURI

	// ShortURIAudit map key is short uri audit entity id attribute
	ShortURIAudit map[string]*_model.ShortURIAudit

	// ShortURIUser map key is short uri user entity id attribute
	ShortURIUser map[string]*_model.ShortURIUser

	shortURIAnchor     sync.RWMutex
	shortURIUserAnchor sync.RWMutex
}

var inMemDB *inMemoryDB

func newInMemoryDB(kind string) (*inMemoryDB, error) {
	if inMemDB != nil {
		return inMemDB, nil
	}

	inMemDB = &inMemoryDB{
		ShortURI:      make(map[string]*_model.ShortURI),
		ShortURIAudit: make(map[string]*_model.ShortURIAudit),
		ShortURIUser:  make(map[string]*_model.ShortURIUser),
		DBKind:        kind,
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

func (db *inMemoryDB) GetShortURIRWMutex() *sync.RWMutex {
	return &db.shortURIAnchor
}

func (db *inMemoryDB) GetShortURIUserRWMutex() *sync.RWMutex {
	return &db.shortURIUserAnchor
}

func (db *inMemoryDB) GetShortURICache() map[string]*_model.ShortURI {
	return db.ShortURI
}

func (db *inMemoryDB) GetShortURIAuditCache() map[string]*_model.ShortURIAudit {
	return db.ShortURIAudit
}

func (db *inMemoryDB) GetShortURIUserCache() map[string]*_model.ShortURIUser {
	return db.ShortURIUser
}

// ========
