package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
)

const (
	getShortURISQL           string = "select id, original_url, key from short_uris where id = $1"
	getShortURIByKeySQL      string = "select id, original_url, key from short_uris where key = $1"
	getShortURIByKeyUserSQL  string = `select su.id, su.original_url, su.key, suu.deleted from short_uris su inner join short_uri_users suu on suu.short_uri_id = su.id and suu.user_id = $2 where su.key = $1`
	createShortURISQL        string = "insert into short_uris(id, original_url, key) values ($1, $2, $3)"
	listShortURIAllByUserSQL string = `select
    s.id,
    s.original_url,
    s.key
from
    short_uris s
    inner join short_uri_users su
        on
            su.user_id = $1
        and su.short_uri_id = s.id`
	listShortURIIdsByKeysSQL string = `select su.id from short_uris su where su.key = any($1)`
)

type shortURIPgRepo struct {
	db       _db.DB
	userRepo ShortURIUserRepository
}

func (pgs *shortURIPgRepo) ListAllByKeys(ctx context.Context, keys []string) ([]*_model.ShortURI, error) {
	//TODO implement me
	panic("implement me")
}

func newShortURIPgRepo(db _db.DB) (*shortURIPgRepo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	userRepo, err := NewShortURIUserRepository(db)
	if err != nil {
		return nil, err
	}

	return &shortURIPgRepo{
		db:       db,
		userRepo: userRepo,
	}, nil
}

func (pgs *shortURIPgRepo) Get(ctx context.Context, id string) (*_model.ShortURI, error) {
	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURISQL, id)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}

	// id, original_url, key
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgs *shortURIPgRepo) GetByKey(ctx context.Context, key string) (*_model.ShortURI, error) {
	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURIByKeySQL, key)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}

	// id, original_url, key, create_user, created, update_user, updated
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgs *shortURIPgRepo) GetByKeyUser(ctx context.Context, userID string, key string) (*_model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	if key == "" {
		return nil, nil
	}

	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURIByKeyUserSQL, key, userID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}
	var deleted = false
	err := row.Scan(&result.ID, &result.OriginalURL, &result.Key, &deleted)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if deleted {
		return nil, _err.NewAppSoftRemovedError("short_uri", nil)
	}

	return &result, nil
}

func (pgs *shortURIPgRepo) Create(ctx context.Context, userID string, entity *_model.ShortURI) (*_model.ShortURI, error) {
	if err := _model.ValidateShortURI(entity); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, _err.NewAppAuthInfoAbsentError("short_uri", nil)
	}

	find, err := pgs.GetByKey(ctx, entity.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		err := pgs.addUser(ctx, nil, find.ID, userID)
		if err != nil && errors.As(err, &_err.AppModelAlreadyExists) {
			return find, err
		} else if err != nil {
			return nil, err
		}
		return find, errors.New("shortURI already exists")
	}

	tx, err := pgs.db.GetDB().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			return
		}

		tx.Commit()
	}()

	stmt, err := tx.PrepareContext(ctx, createShortURISQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)
	stmtSU, err := tx.PrepareContext(ctx, createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmtSU)

	// id, original_url, key
	res, err := pgs.internalCreate(ctx, stmt, entity)
	if err != nil {
		return nil, err
	}
	if err := pgs.addUser(ctx, stmtSU, res.ID, userID); err != nil {
		return nil, err
	}

	return res, nil
}

// BatchCreate is creation a batch data in transaction
func (pgs *shortURIPgRepo) BatchCreate(ctx context.Context, userID string, batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	if len(batch) == 0 {
		return batch, nil
	}

	for _, entity := range batch {
		if err := _model.ValidateShortURI(entity); err != nil {
			return nil, fmt.Errorf("batch validation, invalid entity: [%v] with error [%v]", entity, err)
		}
	}
	if userID == "" {
		return nil, _err.NewAppAuthInfoAbsentError("short uri batch create", nil)
	}

	tx, err := pgs.db.GetDB().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			return
		}

		tx.Commit()
	}()

	stmt, err := tx.Prepare(createShortURISQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)
	stmtSU, err := tx.PrepareContext(ctx, createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmtSU)

	res := make(map[string]*_model.ShortURI)
	for correlation, entity := range batch {
		find, err := pgs.GetByKey(ctx, entity.Key)
		if err != nil {
			return nil, err
		}
		if find != nil {
			if err := pgs.addUser(ctx, stmtSU, find.ID, userID); err != nil {
				return nil, err
			}
			res[correlation] = find

			continue
		}

		saved, err := pgs.internalCreate(ctx, stmt, entity)
		if err != nil {
			return nil, err
		}
		if err := pgs.addUser(ctx, stmtSU, saved.ID, userID); err != nil {
			return nil, err
		}

		res[correlation] = saved
	}

	return res, nil
}

func (pgs *shortURIPgRepo) ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURI, error) {
	res := make([]*_model.ShortURI, 0)
	if userID == "" {
		return res, nil
	}
	rows, err := pgs.db.GetDB().QueryContext(ctx, listShortURIAllByUserSQL, userID)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(rows)
	for rows.Next() {
		var result = _model.ShortURI{
			OriginalURL: &_model.CustomURL{},
		}

		err := rows.Scan(&result.ID, result.OriginalURL, &result.Key)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgs *shortURIPgRepo) Delete(ctx context.Context, ID string, userID string) error {
	return pgs.userRepo.DeleteByUnique(ctx, userID, ID)
}

func (pgs *shortURIPgRepo) BatchDeleteByKeys(ctx context.Context, userID string, keys []string) error {
	if userID == "" || len(keys) == 0 {
		return nil
	}

	ids, err := pgs.listIdsByKeys(ctx, keys)
	if err != nil {
		return err
	}

	return pgs.userRepo.DeleteAllByUnique(ctx, userID, ids)
}

func (pgs *shortURIPgRepo) internalCreate(ctx context.Context, preparedSQL *sql.Stmt, entity *_model.ShortURI) (*_model.ShortURI, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	_, err = preparedSQL.ExecContext(ctx, entity.ID, entity.OriginalURL, entity.Key)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pgs *shortURIPgRepo) addUser(ctx context.Context, stmt *sql.Stmt, id string, userID string) error {
	user, err := _model.NewShortURIUser(id, userID)
	if err != nil {
		return err
	}

	if _, err := pgs.userRepo.CreateStmt(ctx, stmt, user); err != nil {
		return err
	}

	return nil
}

func (pgs *shortURIPgRepo) listIdsByKeys(ctx context.Context, keys []string) ([]string, error) {
	res := make([]string, 0)
	if len(keys) == 0 {
		return res, nil
	}

	rows, err := pgs.db.GetDB().QueryContext(ctx, listShortURIIdsByKeysSQL, keys)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, id)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}
