package db

import (
	"database/sql"
	"sync"

	"github.com/ElfAstAhe/url-shortener/internal/model"
)

type inMemoryDB struct {
	DBKind string

	// ShortURI map key is short uri entity id attribute
	ShortURI map[string]*model.ShortURI

	// ShortURIAudit map key is short uri audit entity id attribute
	ShortURIAudit map[string]*model.ShortURIAudit

	// ShortURIUser map key is short uri user entity id attribute
	ShortURIUser map[string]*model.ShortURIUser

	shortURIAnchor     sync.RWMutex
	shortURIUserAnchor sync.RWMutex
}

var inMemDB *inMemoryDB

func newInMemoryDB(kind string) (*inMemoryDB, error) {
	if inMemDB != nil {
		return inMemDB, nil
	}

	inMemDB = &inMemoryDB{
		ShortURI:      make(map[string]*model.ShortURI),
		ShortURIAudit: make(map[string]*model.ShortURIAudit),
		ShortURIUser:  make(map[string]*model.ShortURIUser),
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

func (db *inMemoryDB) GetShortURICache() map[string]*model.ShortURI {
	return db.ShortURI
}

func (db *inMemoryDB) GetShortURIAuditCache() map[string]*model.ShortURIAudit {
	return db.ShortURIAudit
}

func (db *inMemoryDB) GetShortURIUserCache() map[string]*model.ShortURIUser {
	return db.ShortURIUser
}

// ========
