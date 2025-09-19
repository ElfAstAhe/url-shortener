package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
	"github.com/google/uuid"
)

const (
	getShortURISQL           string = "select id, original_url, key from short_uris where id = $1"
	getShortURIByKeySQL      string = "select id, original_url, key from short_uris where key = $1"
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
)

type shortURIPgRepo struct {
	db       _db.DB
	userRepo ShortURIUserRepository
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

func (pgs *shortURIPgRepo) Create(ctx context.Context, entity *_model.ShortURI) (*_model.ShortURI, error) {
	if err := _model.ValidateShortURI(entity); err != nil {
		return nil, err
	}

	userInfo, err := _auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	find, err := pgs.GetByKey(ctx, entity.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		if err := pgs.addUser(ctx, nil, find.ID, userInfo.UserID); err != nil {
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

	preparedSQL, err := tx.Prepare(createShortURISQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(preparedSQL)
	// id, original_url, key
	res, err := pgs.internalCreate(ctx, preparedSQL, entity)
	if err != nil {
		return nil, err
	}
	if err := pgs.addUser(ctx, tx, res.ID, userInfo.UserID); err != nil {
		return nil, err
	}

	return res, nil
}

// BatchCreate is creation a batch data in transaction
func (pgs *shortURIPgRepo) BatchCreate(ctx context.Context, batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	if len(batch) == 0 {
		return batch, nil
	}

	for _, entity := range batch {
		if err := _model.ValidateShortURI(entity); err != nil {
			return nil, fmt.Errorf("batch validation, invalid entity: [%v] with error [%v]", entity, err)
		}
	}

	userInfo, err := _auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
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

	res := make(map[string]*_model.ShortURI)
	for correlation, entity := range batch {
		find, err := pgs.GetByKey(ctx, entity.Key)
		if err != nil {
			return nil, err
		}
		if find != nil {
			if err := pgs.addUser(ctx, tx, find.ID, userInfo.UserID); err != nil {
				return nil, err
			}
			res[correlation] = find

			continue
		}

		saved, err := pgs.internalCreate(ctx, stmt, entity)
		if err != nil {
			return nil, err
		}
		if err := pgs.addUser(ctx, tx, saved.ID, userInfo.UserID); err != nil {
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

func (pgs *shortURIPgRepo) addUser(ctx context.Context, tx *sql.Tx, id string, userID string) error {
	find, err := pgs.userRepo.GetByUnique(ctx, userID, id)
	if err != nil {
		return err
	}
	if find != nil {
		return nil
	}

	user, err := _model.NewShortURIUser(id, userID)
	if err != nil {
		return err
	}

	if _, err := pgs.userRepo.CreateTran(ctx, tx, user); err != nil {
		return err
	}

	return nil
}
