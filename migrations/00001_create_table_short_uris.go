package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(UpCreateTableShortUris, DownCreateTableShortUris)
}

func UpCreateTableShortUris(ctx context.Context, db *sql.DB) error {
	// ToDo: implementation

	return nil
}

func DownCreateTableShortUris(ctx context.Context, db *sql.DB) error {
	// ToDo: implementation

	return nil
}
