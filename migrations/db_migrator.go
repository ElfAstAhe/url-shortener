package migrations

type DBMigrator interface {
	Initialize() error
	Up() error
}
