package repository

import (
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type shortUriInMemRepo struct {
	Db *_db.InMemoryDb
}

func NewShortUriInMemRepo() ShortUriRepository {
	return &shortUriInMemRepo{
		Db: _db.InMemoryDbInstance,
	}
}

func (r *shortUriInMemRepo) GetById(id string) (*_model.ShortUri, error) {
	for _, value := range r.Db.ShortUri {
		if value.Id == id {
			return value, nil
		}
	}

	return nil, nil
}

func (r *shortUriInMemRepo) GetByKey(key string) (*_model.ShortUri, error) {
	res := r.Db.ShortUri[key]

	return res, nil
}

func (r *shortUriInMemRepo) Create(shortUri *_model.ShortUri) (*_model.ShortUri, error) {
	founded, err := r.GetByKey(shortUri.Key)
	if err != nil {
		return nil, err
	}
	if founded != nil {
		return founded, nil
	}

	r.Db.ShortUri[shortUri.Key] = shortUri

	return shortUri, nil
}
