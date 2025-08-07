package repository

import (
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type shortURIInMemRepo struct {
	DB *_db.InMemoryDB
}

func NewShortURIInMemRepo() ShortURIRepository {
	return &shortURIInMemRepo{
		DB: _db.InMemoryDBInstance,
	}
}

func (r *shortURIInMemRepo) GetByID(id string) (*_model.ShortURI, error) {
	for _, value := range r.DB.ShortURI {
		if value.ID == id {
			return value, nil
		}
	}

	return nil, nil
}

func (r *shortURIInMemRepo) GetByKey(key string) (*_model.ShortURI, error) {
	res := r.DB.ShortURI[key]

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

	r.DB.ShortURI[shortURI.Key] = shortURI

	return shortURI, nil
}
