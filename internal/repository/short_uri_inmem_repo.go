package repository

import (
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type shortURIInMemRepo struct {
	Db *_db.InMemoryDB
}

func NewShortUriInMemRepo() ShortURIRepository {
	return &shortURIInMemRepo{
		Db: _db.InMemoryDBInstance,
	}
}

func (r *shortURIInMemRepo) GetById(id string) (*_model.ShortURI, error) {
	for _, value := range r.Db.ShortURI {
		if value.Id == id {
			return value, nil
		}
	}

	return nil, nil
}

func (r *shortURIInMemRepo) GetByKey(key string) (*_model.ShortURI, error) {
	res := r.Db.ShortURI[key]

	return res, nil
}

func (r *shortURIInMemRepo) Create(shortUri *_model.ShortURI) (*_model.ShortURI, error) {
	founded, err := r.GetByKey(shortUri.Key)
	if err != nil {
		return nil, err
	}
	if founded != nil {
		return founded, nil
	}

	r.Db.ShortURI[shortUri.Key] = shortUri

	return shortUri, nil
}
