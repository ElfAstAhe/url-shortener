package repository

import (
	"context"
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/google/uuid"
)

type shortURIUserInMemRepo struct {
	Cache _db.InMemoryCache
}

func newShortURIUserInMemRepo(db _db.DB) (*shortURIUserInMemRepo, error) {
	if cache, ok := db.(_db.InMemoryCache); ok {
		return &shortURIUserInMemRepo{
			Cache: cache,
		}, nil
	}

	return nil, errors.New("db param does not implement InMemoryCache")
}

func (imsu *shortURIUserInMemRepo) Get(ctx context.Context, ID string) (*_model.ShortURIUser, error) {
	if ID == "" {
		return nil, nil
	}

	return imsu.Cache.GetShortURIUserCache()[ID], nil
}

func (imsu *shortURIUserInMemRepo) GetByUnique(ctx context.Context, userID string, shortURIID string) (*_model.ShortURIUser, error) {
	if userID == "" || shortURIID == "" {
		return nil, nil
	}

	for _, entity := range imsu.Cache.GetShortURIUserCache() {
		if entity.UserID == userID && entity.ShortURIID == shortURIID {
			return entity, nil
		}
	}

	return nil, nil
}

func (imsu *shortURIUserInMemRepo) ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURIUser, error) {
	res := make([]*_model.ShortURIUser, 0)
	if userID == "" || len(imsu.Cache.GetShortURIUserCache()) == 0 {
		return res, nil
	}

	for _, entity := range imsu.Cache.GetShortURIUserCache() {
		if entity.UserID == userID {
			res = append(res, entity)
		}
	}

	return res, nil
}

func (imsu *shortURIUserInMemRepo) ListAllByShortURI(ctx context.Context, shortURIID string) ([]*_model.ShortURIUser, error) {
	res := make([]*_model.ShortURIUser, 0)
	if shortURIID == "" || len(imsu.Cache.GetShortURIUserCache()) == 0 {
		return res, nil
	}

	for _, entity := range imsu.Cache.GetShortURIUserCache() {
		if entity.ShortURIID == shortURIID {
			res = append(res, entity)
		}
	}

	return res, nil
}

func (imsu *shortURIUserInMemRepo) CreateTran(ctx context.Context, tx *sql.Tx, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	return imsu.Create(ctx, entity)
}

func (imsu *shortURIUserInMemRepo) Create(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	entity.ID = newID.String()
	imsu.Cache.GetShortURIUserCache()[entity.ID] = entity

	return entity, nil
}

func (imsu *shortURIUserInMemRepo) Change(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	imsu.Cache.GetShortURIUserCache()[entity.ID] = entity

	return entity, nil
}
