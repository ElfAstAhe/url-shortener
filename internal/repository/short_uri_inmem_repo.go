package repository

import (
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/google/uuid"
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

func (imr *shortURIInMemRepo) Get(id string) (*_model.ShortURI, error) {
	for _, value := range imr.Cache.GetShortURICache() {
		if value.ID == id {
			return value, nil
		}
	}

	return nil, nil
}

func (imr *shortURIInMemRepo) GetByKey(key string) (*_model.ShortURI, error) {
	res := imr.Cache.GetShortURICache()[key]

	return res, nil
}

func (imr *shortURIInMemRepo) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	if shortURI == nil {
		return nil, errors.New("shortURI is nil")
	}
	if shortURI.Key == "" {
		return nil, errors.New("shortURI Key is empty")
	}

	founded, err := imr.GetByKey(shortURI.Key)
	if err != nil {
		return nil, err
	}
	if founded != nil {
		return founded, errors.New("short URI already exists")
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	shortURI.ID = newID.String()

	imr.Cache.GetShortURICache()[shortURI.Key] = shortURI

	return shortURI, nil
}

func (imr *shortURIInMemRepo) BatchCreate(batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	if len(batch) == 0 {
		return batch, nil
	}

	res := make(map[string]*_model.ShortURI)
	for correlation, item := range batch {
		// searching
		find, err := imr.GetByKey(item.Key)
		if err != nil {
			return nil, err
		}
		// founded
		if find != nil {
			res[correlation] = find

			continue
		}
		// new one
		data, err := imr.Create(item)
		if err != nil {
			return nil, err
		}
		res[correlation] = data
	}

	return res, nil
}

func (imr *shortURIInMemRepo) ListAllByUser(userID string) ([]*_model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	// all shorten ids by user
	ids, err := imr.listIDsByUser(userID)
	if err != nil {
		return nil, err
	}

	return imr.listAllByIDs(ids)
}

func (imr *shortURIInMemRepo) listIDsByUser(userID string) ([]string, error) {
	if userID == "" {
		return nil, nil
	}
	res := make([]string, 0)
	for _, value := range imr.Cache.GetShortURIUserCache() {
		if value.UserID == userID {
			res = append(res, value.ShortURIID)
		}
	}

	return res, nil
}

func (imr *shortURIInMemRepo) listAllByIDs(ids []string) ([]*_model.ShortURI, error) {
	if len(ids) == 0 {
		return []*_model.ShortURI{}, nil
	}
	res := make([]*_model.ShortURI, 0)
	for _, id := range ids {
		find, err := imr.Get(id)
		if err != nil {
			return nil, err
		}
		if find != nil {
			res = append(res, find)
		}
	}

	return res, nil
}
