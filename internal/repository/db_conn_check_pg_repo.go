package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ElfAstAhe/url-shortener/internal/config/db"
)

type DBConnCheckPgRepo struct {
	DB db.DB
}

func newDBConnCheckPgRepo(appDb db.DB) (*DBConnCheckPgRepo, error) {
	if appDb == nil {
		return nil, errors.New("db is nil")
	}

	return &DBConnCheckPgRepo{
		DB: appDb,
	}, nil
}

// Closer

func (pgR *DBConnCheckPgRepo) Close() error {
	return db.CloseDB(pgR.DB)
}

// ========

// DBConn

func (pgR *DBConnCheckPgRepo) CheckDBConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return pgR.DB.GetDB().PingContext(ctx)
}

// ========
