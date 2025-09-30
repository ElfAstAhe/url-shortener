package repository

import (
	"context"
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
)

type shortURIUserPgRepo struct {
	db _db.DB
}

const (
	getShortURIUserSQL                 string = `select id, short_uri_id, user_id, deleted from short_uri_users where id = $1`
	getShortURIUserByUniqueSQL         string = `select id, short_uri_id, user_id, deleted from short_uri_users where user_id = $1 and short_uri_id = $2`
	listShortURIUserAllByUserSQL       string = `select id, short_uri_id, user_id, deleted from short_uri_users where user_id = $1`
	listShortURIUserAllByShortURISQL   string = `select id, short_uri_id, user_id, deleted from short_uri_users where short_uri_id = $1`
	createShortURIUserSQL              string = `insert into short_uri_users(id, short_uri_id, user_id, deleted) values($1, $2, $3, $4)`
	changeShortURIUserSQL              string = `update short_uri_users set short_uri_id=$2, user_id=$3, deleted = $4 where id = $1`
	deleteShortURIUserSQL              string = `update short_uri_users set deleted = true where id = $1`
	deleteShortURIUserByUniqueSQL      string = `update short_uri_users set deleted = true where user_id = $1 and short_uri_id = $2`
	deleteAllShortURIUserByUniqueSQL   string = `update short_uri_users set deleted = true where user_id = $1 and short_uri_id = any($2)`
	deleteAllShortURIUserByUserSQL     string = `update short_uri_users set deleted = true where user_id = $1`
	deleteAllShortURIUserByShortURISQL string = `update short_uri_users set deleted = true where short_uri_id = $1`
	removeShortURIUserSQL              string = `delete from short_uri_users where id = $1`
	removeShortURIUserByUniqueSQL      string = `delete from short_uri_users where user_id = $1 and short_uri_id = $2`
	removeAllShortURIUserByUniqueSQL   string = `delete from short_uri_users where user_id = $1 and short_uri_id = any($2)`
	removeAllShortURIUserByUserSQL     string = `delete from short_uri_users where user_id = $1`
	removeAllShortURIUserByShortURISQL string = `delete from short_uri_users where short_uri_id = $1`
)

func newShortURIUserPgRepo(db _db.DB) (*shortURIUserPgRepo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &shortURIUserPgRepo{db: db}, nil
}

func (pgsu *shortURIUserPgRepo) Get(ctx context.Context, ID string) (*_model.ShortURIUser, error) {
	row := pgsu.db.GetDB().QueryRow(getShortURIUserSQL, ID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURIUser{}
	// id, short_uri_id, user_id, deleted
	err := row.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgsu *shortURIUserPgRepo) GetByUnique(ctx context.Context, userID string, shortURIID string) (*_model.ShortURIUser, error) {
	row := pgsu.db.GetDB().QueryRow(getShortURIUserByUniqueSQL, userID, shortURIID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURIUser{}
	// id, short_uri_id, user_id, deleted
	err := row.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgsu *shortURIUserPgRepo) ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURIUser, error) {
	res := make([]*_model.ShortURIUser, 0)
	rows, err := pgsu.db.GetDB().Query(listShortURIUserAllByUserSQL, userID)
	if err != nil {
		return res, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		var result = _model.ShortURIUser{}
		// id, short_uri_id, user_id, deleted
		err := rows.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgsu *shortURIUserPgRepo) ListAllByShortURI(ctx context.Context, shortURIID string) ([]*_model.ShortURIUser, error) {
	res := make([]*_model.ShortURIUser, 0)
	rows, err := pgsu.db.GetDB().Query(listShortURIUserAllByShortURISQL, shortURIID)
	if err != nil {
		return res, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		var result = _model.ShortURIUser{}
		// id, short_uri_id, user_id, deleted
		err := rows.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgsu *shortURIUserPgRepo) Create(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)

	res, err := pgsu.CreateStmt(ctx, stmt, entity)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (pgsu *shortURIUserPgRepo) CreateStmt(ctx context.Context, stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	find, err := pgsu.GetByUnique(ctx, entity.UserID, entity.ShortURIID)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, _err.NewAppModelAlreadyExistsError(entity.ID, "short_uri_user")
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	_, err = stmt.ExecContext(ctx, entity.ID, entity.ShortURIID, entity.UserID, entity.Deleted)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pgsu *shortURIUserPgRepo) Change(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, _err.NewAppModelValidationError("short_uri_user", err)
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, changeShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.ChangeStmt(ctx, stmt, entity)
}

func (pgsu *shortURIUserPgRepo) ChangeStmt(ctx context.Context, stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, _err.NewAppModelValidationError("short_uri_user", err)
	}

	_, err := stmt.ExecContext(ctx, entity.ID, entity.ShortURIID, entity.UserID, entity.Deleted)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pgsu *shortURIUserPgRepo) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteShortURIUserSQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.DeleteStmt(ctx, stmt, ID)
}

func (pgsu *shortURIUserPgRepo) DeleteStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	if ID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) DeleteByUnique(ctx context.Context, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.DeleteByUniqueStmt(ctx, stmt, userID, shortURIID)
}

func (pgsu *shortURIUserPgRepo) DeleteByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) DeleteAllByUnique(ctx context.Context, userID string, shortURIIDs []string) error {
	if userID == "" || len(shortURIIDs) == 0 {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteAllShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}

	return pgsu.DeleteAllByUniqueStmt(ctx, stmt, userID, shortURIIDs)
}

func (pgsu *shortURIUserPgRepo) DeleteAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIDs []string) error {
	if userID == "" || len(shortURIIDs) == 0 {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIIDs)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) DeleteAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteAllShortURIUserByUserSQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.DeleteAllByUserStmt(ctx, stmt, userID)
}

func (pgsu *shortURIUserPgRepo) DeleteAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	if userID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) DeleteAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteAllShortURIUserByShortURISQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.DeleteAllByShortURIStmt(ctx, stmt, shortURIID)
}

func (pgsu *shortURIUserPgRepo) DeleteAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, shortURIID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) Remove(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeShortURIUserSQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.RemoveStmt(ctx, stmt, ID)
}

func (pgsu *shortURIUserPgRepo) RemoveStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	if ID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) RemoveByUnique(ctx context.Context, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	stmp, err := pgsu.db.GetDB().PrepareContext(ctx, removeShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmp)

	return pgsu.RemoveByUniqueStmt(ctx, stmp, userID, shortURIID)
}

func (pgsu *shortURIUserPgRepo) RemoveByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) RemoveAllByUnique(ctx context.Context, userID string, shortURIIDs []string) error {
	if userID == "" || len(shortURIIDs) == 0 {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeAllShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}

	return pgsu.RemoveAllByUniqueStmt(ctx, stmt, userID, shortURIIDs)
}

func (pgsu *shortURIUserPgRepo) RemoveAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIds []string) error {
	if userID == "" || len(shortURIIds) == 0 {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIIds)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) RemoveAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeAllShortURIUserByUserSQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.RemoveAllByUserStmt(ctx, stmt, userID)
}

func (pgsu *shortURIUserPgRepo) RemoveAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	if userID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *shortURIUserPgRepo) RemoveAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeAllShortURIUserByShortURISQL)
	if err != nil {
		return err
	}
	defer _utl.CloseOnly(stmt)

	return pgsu.RemoveAllByShortURIStmt(ctx, stmt, shortURIID)
}

func (pgsu *shortURIUserPgRepo) RemoveAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, shortURIID)
	if err != nil {
		return err
	}

	return nil
}
