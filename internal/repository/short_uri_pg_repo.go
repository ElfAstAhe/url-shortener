package repository

import (
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
	"github.com/google/uuid"
)

const (
	getShortURISQL           string = "select id, original_url, key, create_user, created, update_user, updated from short_uris where id = $1"
	getShortURIByKeySQL      string = "select id, original_url, key, create_user, created, update_user, updated from short_uris where key = $1"
	createShortURISQL        string = "insert into short_uris(id, original_url, key, create_user, created, update_user, updated) values ($1, $2, $3, $4, $5, $6, $7)"
	listShortURIAllByUserSQL string = `select
    s.id,
    s.original_url,
    s.key,
    s.create_user,
    s.created,
    s.update_user,
    s.updated
from
    short_uris s
    inner join short_uri_users su
        on
            su.user_id = $1
        and su.short_uri_id = s.id`
	createUserSQL string = `insert into short_uri_users(id, user_id, short_uri_id) values ($1, $2, $3)`
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

func (pgs *shortURIPgRepo) Get(id string) (*_model.ShortURI, error) {
	row := pgs.db.GetDB().QueryRow(getShortURISQL, id)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}

	// id, original_url, key, create_user, created, update_user, updated
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key, &result.CreateUser, &result.Created, &result.UpdateUser, &result.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgs *shortURIPgRepo) GetByKey(key string) (*_model.ShortURI, error) {
	row := pgs.db.GetDB().QueryRow(getShortURIByKeySQL, key)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}

	// id, original_url, key, create_user, created, update_user, updated
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key, &result.CreateUser, &result.Created, &result.UpdateUser, &result.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgs *shortURIPgRepo) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	if shortURI == nil {
		return nil, errors.New("shortURI is nil")
	}
	if shortURI.Key == "" {
		return nil, errors.New("shortURI Key is empty")
	}

	find, err := pgs.GetByKey(shortURI.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, errors.New("shortURI already exists")
	}

	preparedSQL, err := pgs.db.GetDB().Prepare(createShortURISQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(preparedSQL)
	// id, original_url, key, create_user, created, update_user, updated
	res, err := internalCreateShortURI(preparedSQL, shortURI)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// BatchCreate is creation a batch data in transaction
func (pgs *shortURIPgRepo) BatchCreate(batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	if len(batch) == 0 {
		return batch, nil
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
	for correlation, shortURL := range batch {
		find, err := pgs.GetByKey(shortURL.Key)
		if err != nil {
			return nil, err
		}
		if find != nil {
			res[correlation] = find

			continue
		}
		saved, err := internalCreateShortURI(stmt, shortURL)
		if err != nil {
			return nil, err
		}
		res[correlation] = saved
	}

	return res, nil
}

func (pgs *shortURIPgRepo) ListAllByUser(userID string) ([]*_model.ShortURI, error) {
	rows, err := pgs.db.GetDB().Query(listShortURIAllByUserSQL, userID)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(rows)
	var res = make([]*_model.ShortURI, 0)
	for rows.Next() {
		var result = _model.ShortURI{
			OriginalURL: &_model.CustomURL{},
		}

		err := rows.Scan(&result.ID, result.OriginalURL, &result.Key, &result.CreateUser, &result.Created, &result.UpdateUser, &result.Updated)
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

func internalCreateShortURI(preparedSQL *sql.Stmt, shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	shortURI.ID = newID.String()

	_, err = preparedSQL.Exec(shortURI.ID, shortURI.OriginalURL, shortURI.Key, shortURI.CreateUser, shortURI.Created, shortURI.UpdateUser, shortURI.Updated)
	if err != nil {
		return nil, err
	}

	return shortURI, nil
}
