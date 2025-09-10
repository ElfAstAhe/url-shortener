package repository

import (
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/google/uuid"
)

const (
	getSql      string = "select id, original_url, key, create_user, created, update_user, updated from short_uris where id = $1"
	getByKeySql string = "select id, original_url, key, create_user, created, update_user, updated from short_uris where key = $1"
	createSql   string = "insert into short_uris(id, original_url, key, create_user, created, update_user, updated) values ($1, $2, $3, $4, $5, $6, $7)"
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
	row := pgr.db.GetDB().QueryRow(getSql, id)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}

	// id, original_url, key, create_user, created, update_user, updated
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key, &result.CreateUser, &result.Created, &result.UpdateUser, &result.Updated)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgr *shortURIPgRepo) GetByKey(key string) (*_model.ShortURI, error) {
	row := pgr.db.GetDB().QueryRow(getByKeySql, key)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURI{
		OriginalURL: &_model.CustomURL{},
	}

	// id, original_url, key, create_user, created, update_user, updated
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key, &result.CreateUser, &result.Created, &result.UpdateUser, &result.Updated)
	if err != nil {
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
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	shortURI.ID = newID.String()

	// id, original_url, key, create_user, created, update_user, updated
	_, err = pgr.db.GetDB().Exec(createSql, shortURI.ID, shortURI.OriginalURL, shortURI.Key, shortURI.CreateUser, shortURI.Created, shortURI.UpdateUser, shortURI.Updated)
	if err != nil {
		return nil, err
	}

	return pgr.GetByKey(shortURI.Key)
}
