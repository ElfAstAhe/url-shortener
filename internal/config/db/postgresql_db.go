package db

import (
	"database/sql"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresqlDB struct {
	DB     *sql.DB
	DBKind string
	Dsn    string
}

var pgDB *postgresqlDB

func newPostgresqlDB(kind string, dsn string) (*postgresqlDB, error) {
	if pgDB != nil {
		return pgDB, nil
	}

	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pgDB = &postgresqlDB{
		DB:     pg,
		DBKind: kind,
		Dsn:    dsn,
	}

	return pgDB, nil
}

func NewPGIter10Gap(dsn string) (DB, error) {
	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &postgresqlDB{
		DB:     pg,
		DBKind: _cfg.DBKindPostgres,
		Dsn:    dsn,
	}, nil
}

// Closer

func (pDB *postgresqlDB) Close() error {
	return pDB.DB.Close()
}

// =============

// DB

func (pDB *postgresqlDB) GetDB() *sql.DB {
	return pDB.DB
}

func (pDB *postgresqlDB) GetDBKind() string {
	return pDB.DBKind
}

func (pDB *postgresqlDB) GetDsn() string {
	return pDB.Dsn
}

// =============
