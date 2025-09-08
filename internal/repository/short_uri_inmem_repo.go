package repository

import (
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type shortURIInMemRepo struct {
	Cache _db.InMemoryCache
}

func newShortURIInMemRepo(db _db.DB) (ShortURIRepository, error) {
	if cache, ok := db.(_db.InMemoryCache); ok {
		return &shortURIInMemRepo{
			Cache: cache,
		}, nil
	}

	return nil, errors.New("db param does not implement InMemoryCache")
}

func (r *shortURIInMemRepo) Get(id string) (*_model.ShortURI, error) {
	for _, value := range r.Cache.GetShortURICache() {
		if value.ID == id {
			return value, nil
		}
	}

	return nil, nil
}

func (r *shortURIInMemRepo) GetByKey(key string) (*_model.ShortURI, error) {
	res := r.Cache.GetShortURICache()[key]

	return res, nil
}

func (r *shortURIInMemRepo) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	founded, err := r.GetByKey(shortURI.Key)
	if err != nil {
		return nil, err
	}
	if founded != nil {
		return founded, nil
	}

	r.Cache.GetShortURICache()[shortURI.Key] = shortURI

	return shortURI, nil
}
