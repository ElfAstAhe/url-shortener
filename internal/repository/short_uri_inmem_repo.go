package repository

import (
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/google/uuid"
)

type shortURIInMemRepo struct {
	Cache    _db.InMemoryCache
	userRepo ShortURIUserRepository
}

func newShortURIInMemRepo(db _db.DB) (*shortURIInMemRepo, error) {
	if cache, ok := db.(_db.InMemoryCache); ok {
		userRepo, err := NewShortURIUserRepository(db)
		if err != nil {
			return nil, err
		}

		return &shortURIInMemRepo{
			Cache:    cache,
			userRepo: userRepo,
		}, nil
	}

	return nil, errors.New("db param does not implement InMemoryCache")
}

func (ims *shortURIInMemRepo) Get(id string) (*_model.ShortURI, error) {
	for _, value := range ims.Cache.GetShortURICache() {
		if value.ID == id {
			return value, nil
		}
	}

	return nil, nil
}

func (ims *shortURIInMemRepo) GetByKey(key string) (*_model.ShortURI, error) {
	res := ims.Cache.GetShortURICache()[key]

	return res, nil
}

func (ims *shortURIInMemRepo) Create(shortURI *_model.ShortURI) (*_model.ShortURI, error) {
	if shortURI == nil {
		return nil, errors.New("shortURI is nil")
	}
	if shortURI.Key == "" {
		return nil, errors.New("shortURI Key is empty")
	}

	find, err := ims.GetByKey(shortURI.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		if err := ims.addUser(find.ID, shortURI.CreateUser); err != nil {
			return nil, err
		}

		return find, errors.New("short URI already exists")
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	shortURI.ID = newID.String()

	ims.Cache.GetShortURICache()[shortURI.Key] = shortURI

	if err := ims.addUser(shortURI.ID, shortURI.CreateUser); err != nil {
		return nil, err
	}

	return shortURI, nil
}

func (ims *shortURIInMemRepo) BatchCreate(batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	if len(batch) == 0 {
		return batch, nil
	}

	res := make(map[string]*_model.ShortURI)
	for correlation, item := range batch {
		// searching
		find, err := ims.GetByKey(item.Key)
		if err != nil {
			return nil, err
		}
		// founded
		if find != nil {
			if err := ims.addUser(find.ID, item.CreateUser); err != nil {
				return nil, err
			}

			res[correlation] = find

			continue
		}
		// new one
		data, err := ims.Create(item)
		if err != nil {
			return nil, err
		}
		res[correlation] = data
	}

	return res, nil
}

func (ims *shortURIInMemRepo) ListAllByUser(userID string) ([]*_model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	// all shorten ids by user
	ids, err := ims.listIDsByUser(userID)
	if err != nil {
		return nil, err
	}

	return ims.listAllByIDs(ids)
}

func (ims *shortURIInMemRepo) listIDsByUser(userID string) ([]string, error) {
	if userID == "" {
		return nil, nil
	}
	res := make([]string, 0)
	for _, value := range ims.Cache.GetShortURIUserCache() {
		if value.UserID == userID {
			res = append(res, value.ShortURIID)
		}
	}

	return res, nil
}

func (ims *shortURIInMemRepo) listAllByIDs(ids []string) ([]*_model.ShortURI, error) {
	if len(ids) == 0 {
		return []*_model.ShortURI{}, nil
	}
	res := make([]*_model.ShortURI, 0)
	for _, id := range ids {
		find, err := ims.Get(id)
		if err != nil {
			return nil, err
		}
		if find != nil {
			res = append(res, find)
		}
	}

	return res, nil
}

func (ims *shortURIInMemRepo) addUser(ID string, userID string) error {
	find, err := ims.userRepo.GetByUnique(userID, ID)
	if err != nil {
		return err
	}
	if find != nil {
		return nil
	}

	entity, err := _model.NewShortURIUser(ID, userID)
	if err != nil {
		return err
	}

	res, err := ims.userRepo.Create(entity)

	if err != nil && res == nil {
		return nil
	}

	return nil
}
