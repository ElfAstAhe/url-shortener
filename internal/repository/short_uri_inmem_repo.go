package repository

import (
	"github.com/ElfAstAhe/url-shortener/internal/config/db"
	"github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/ElfAstAhe/url-shortener/pkg/errors"
)

type shortUriInMemRepo struct {
	Db db.InMemoryDb
}

func NewShortUriInMemRepo() ShortUriRepository {
	return &shortUriInMemRepo{}
}

func (r *shortUriInMemRepo) GetById(id string) (*model.ShortUri, error) {
	for key, value := range r.Db.ShortUri {
		if key == id {
			return value, nil
		}
	}

	return nil, errors.NewNotFoundError(id)
}

func (r *shortUriInMemRepo) GetByKey(key string) (*model.ShortUri, error) {
	res := r.Db.ShortUri[key]
	if res == nil {
		return nil, errors.NewNotFoundError(key)
	}

	return res, nil
}

func (r *shortUriInMemRepo) Create(shortUri *model.ShortUri) (*model.ShortUri, error) {
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
