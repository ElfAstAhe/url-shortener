package repository

import (
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type shortURIPgRepo struct {
	db _db.DB
}

func newShortURIPgRepo(db _db.DB) (ShortURIRepository, error) {
	return &shortURIPgRepo{db: db}, nil
}

func (s *shortURIPgRepo) Get(id string) (*_model.ShortURI, error) {
	return nil, errors.New("not implemented")
}

func (s *shortURIPgRepo) GetByKey(key string) (*_model.ShortURI, error) {
	return nil, errors.New("not implemented")
}

func (s *shortURIPgRepo) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	return nil, errors.New("not implemented")
}
