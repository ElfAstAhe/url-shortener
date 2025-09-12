package repository

import (
	"context"
	"errors"
	"time"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
)

type DBConnCheckPgRepo struct {
	DB _db.DB
}

func newDBConnCheckPgRepo(db _db.DB) (*DBConnCheckPgRepo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &DBConnCheckPgRepo{
		DB: db,
	}, nil
}

// Closer

func (pgR *DBConnCheckPgRepo) Close() error {
	return _db.CloseDB(pgR.DB)
}

// ========

// DBConn

func (pgR *DBConnCheckPgRepo) CheckDBConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return pgR.DB.GetDB().PingContext(ctx)
}

// ========
