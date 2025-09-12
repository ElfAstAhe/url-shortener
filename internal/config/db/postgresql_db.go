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

func newPostgresqlDB(kind string, dsn string) (*postgresqlDB, error) {
	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &postgresqlDB{
		DB:     pg,
		DBKind: kind,
		Dsn:    dsn,
	}, nil
}

func NewPGIter10Gap(dsn string) (DB, error) {
	return newPostgresqlDB(_cfg.DBKindPostgres, dsn)
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
