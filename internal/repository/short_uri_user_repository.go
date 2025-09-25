package repository

import (
	"context"
	"database/sql"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type ShortURIUserRepository interface {
	Get(ctx context.Context, ID string) (*_model.ShortURIUser, error)
	GetByUnique(ctx context.Context, userID string, shortURIID string) (*_model.ShortURIUser, error)
	ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURIUser, error)
	ListAllByShortURI(ctx context.Context, shortURIID string) ([]*_model.ShortURIUser, error)
	// Create is create new record
	Create(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	// Change is change record attributes
	Change(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	// Delete is record soft remove
	Delete(ctx context.Context, ID string) error
	// DeleteByUnique is record soft remove by unique key params
	DeleteByUnique(ctx context.Context, userID string, shortURIID string) error
	// DeleteAllByUnique is massive records soft delete by unique
	DeleteAllByUnique(ctx context.Context, userID string, shortURIIds []string) error
	// DeleteAllByUser is massive records soft removals by user
	DeleteAllByUser(ctx context.Context, userID string) error
	// DeleteAllByShortURI is massive records soft removals by short uri
	DeleteAllByShortURI(ctx context.Context, shortURIID string) error
	// Remove is physical record removal
	Remove(ctx context.Context, ID string) error
	// RemoveByUnique is physical record removal by unique params
	RemoveByUnique(ctx context.Context, userID string, shortURIID string) error
	// RemoveAllByUser is physical all records removal by user
	RemoveAllByUser(ctx context.Context, userID string) error
	// RemoveAllByShortURI is physical all records removal by short uri
	RemoveAllByShortURI(ctx context.Context, shortURIID string) error

	CreateStmt(ctx context.Context, stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	ChangeStmt(ctx context.Context, stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error)
	DeleteStmt(ctx context.Context, stmt *sql.Stmt, ID string) error
	DeleteByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error
	DeleteAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIds []string) error
	DeleteAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error
	DeleteAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, userID string) error
	RemoveStmt(ctx context.Context, stmt *sql.Stmt, ID string) error
	RemoveByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error
	RemoveAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error
	RemoveAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error
}

func NewShortURIUserRepository(db _db.DB) (ShortURIUserRepository, error) {
	if db != nil && db.GetDBKind() == _cfg.DBKindPostgres {
		return newShortURIUserPgRepo(db)
	}

	return newShortURIUserInMemRepo(db)
}
