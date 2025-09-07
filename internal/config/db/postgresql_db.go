package db

import (
	"database/sql"
	"fmt"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
)

type postgresqlDB struct {
	DB     *sql.DB
	DBKind string
	Config *_cfg.DBConfig
}

func newPostgresqlDB(kind string, config *_cfg.DBConfig) (*postgresqlDB, error) {
	pg, err := sql.Open("pgx", buildDSN(config))
	if err != nil {
		return nil, err
	}

	return &postgresqlDB{
		DB:     pg,
		DBKind: kind,
		Config: config,
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

func (pDB *postgresqlDB) GetConfig() *_cfg.DBConfig {
	return pDB.Config
}

// =============

func buildDSN(config *_cfg.DBConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Database)
}
