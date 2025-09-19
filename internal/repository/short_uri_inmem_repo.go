package repository

import (
	"context"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
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

func (ims *shortURIInMemRepo) Get(ctx context.Context, id string) (*_model.ShortURI, error) {
	res := ims.Cache.GetShortURICache()[id]

	return res, nil
}

func (ims *shortURIInMemRepo) GetByKey(ctx context.Context, key string) (*_model.ShortURI, error) {
	for _, value := range ims.Cache.GetShortURICache() {
		if value.Key == key {
			return value, nil
		}
	}

	return nil, nil
}

func (ims *shortURIInMemRepo) Create(ctx context.Context, entity *_model.ShortURI) (*_model.ShortURI, error) {
	if err := _model.ValidateShortURI(entity); err != nil {
		return nil, err
	}

	userInfo, err := _auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	find, err := ims.GetByKey(ctx, entity.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		if err := ims.addUser(ctx, find.ID, userInfo.UserID); err != nil {
			return nil, err
		}

		return find, errors.New("short URI already exists")
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	ims.Cache.GetShortURICache()[entity.ID] = entity

	if err := ims.addUser(ctx, entity.ID, userInfo.UserID); err != nil {
		return nil, err
	}

	return entity, nil
}

func (ims *shortURIInMemRepo) BatchCreate(ctx context.Context, batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	res := make(map[string]*_model.ShortURI)
	if len(batch) == 0 {
		return res, nil
	}

	if _, err := _auth.UserInfoFromContext(ctx); err != nil {
		return nil, err
	}

	for correlation, item := range batch {
		data, err := ims.Create(ctx, item)
		if err != nil && data == nil {
			return nil, err
		}
		res[correlation] = data
	}

	return res, nil
}

func (ims *shortURIInMemRepo) ListAllByUser(ctx context.Context, userID string) ([]*_model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	// all shorten ids by user
	entityUserLinks, err := ims.userRepo.ListAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return ims.listAllByLinks(ctx, entityUserLinks)
}

func (ims *shortURIInMemRepo) listAllByLinks(ctx context.Context, userLinks []*_model.ShortURIUser) ([]*_model.ShortURI, error) {
	res := make([]*_model.ShortURI, 0)
	if len(userLinks) == 0 {
		return res, nil
	}
	for _, userLink := range userLinks {
		find, err := ims.Get(ctx, userLink.ShortURIID)
		if err != nil {
			return nil, err
		}
		if find != nil {
			res = append(res, find)
		}
	}

	return res, nil
}

func (ims *shortURIInMemRepo) addUser(ctx context.Context, ID string, userID string) error {
	find, err := ims.userRepo.GetByUnique(ctx, userID, ID)
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

	res, err := ims.userRepo.Create(ctx, entity)

	if err != nil && res == nil {
		return err
	}

	return nil
}
