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
	getSQL      string = "select id, original_url, key, create_user, created, update_user, updated from short_uris where id = $1"
	getByKeySQL string = "select id, original_url, key, create_user, created, update_user, updated from short_uris where key = $1"
	createSQL   string = "insert into short_uris(id, original_url, key, create_user, created, update_user, updated) values ($1, $2, $3, $4, $5, $6, $7)"
)

type shortURIPgRepo struct {
	db _db.DB
}

func newShortURIPgRepo(db _db.DB) (ShortURIRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &shortURIPgRepo{db: db}, nil
}

func (pgr *shortURIPgRepo) Get(id string) (*_model.ShortURI, error) {
	row := pgr.db.GetDB().QueryRow(getSQL, id)
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

func (pgr *shortURIPgRepo) GetByKey(key string) (*_model.ShortURI, error) {
	row := pgr.db.GetDB().QueryRow(getByKeySQL, key)
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

func (pgr *shortURIPgRepo) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	if shortURI == nil {
		return nil, errors.New("shortURI is nil")
	}
	if shortURI.Key == "" {
		return nil, errors.New("shortURI Key is empty")
	}
	preparedSQL, err := pgr.db.GetDB().Prepare(createSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(preparedSQL)
	// id, original_url, key, create_user, created, update_user, updated
	res, err := internalCreate(preparedSQL, shortURI)
	if err != nil {
		return nil, err
	}

	return pgr.Get(res.ID)
}

// BatchCreate is creation a batch data in transaction
func (pgr *shortURIPgRepo) BatchCreate(batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	if batch == nil || len(batch) == 0 {
		return batch, nil
	}

	tx, err := pgr.db.GetDB().Begin()
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

	stmt, err := tx.Prepare(createSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)

	res := make(map[string]*_model.ShortURI)
	for correlation, shortURL := range batch {
		find, err := pgr.GetByKey(shortURL.Key)
		if err != nil {
			return nil, err
		}
		if find != nil {
			res[correlation] = find

			continue
		}
		saved, err := internalCreate(stmt, shortURL)
		if err != nil {
			return nil, err
		}
		res[correlation] = saved
	}

	return res, nil
}

func internalCreate(preparedSQL *sql.Stmt, shortURI *_model.ShortURI) (*_model.ShortURI, error) {
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
