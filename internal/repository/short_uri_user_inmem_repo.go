package repository

import (
	"context"
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
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

	imsu.Cache.GetRWMutex().RLock()
	defer imsu.Cache.GetRWMutex().RUnlock()
	return imsu.Cache.GetShortURIUserCache()[ID], nil
}

func (imsu *shortURIUserInMemRepo) GetByUnique(ctx context.Context, userID string, shortURIID string) (*_model.ShortURIUser, error) {
	if userID == "" || shortURIID == "" {
		return nil, nil
	}

	imsu.Cache.GetRWMutex().RLock()
	defer imsu.Cache.GetRWMutex().RUnlock()
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

	imsu.Cache.GetRWMutex().RLock()
	defer imsu.Cache.GetRWMutex().RUnlock()
	for _, entity := range imsu.Cache.GetShortURIUserCache() {
		if entity.UserID == userID && !entity.Deleted {
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

	imsu.Cache.GetRWMutex().RLock()
	defer imsu.Cache.GetRWMutex().RUnlock()
	for _, entity := range imsu.Cache.GetShortURIUserCache() {
		if entity.ShortURIID == shortURIID {
			res = append(res, entity)
		}
	}

	return res, nil
}

func (imsu *shortURIUserInMemRepo) Create(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, _err.NewAppModelValidationError("short_uri_user", err)
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	entity.ID = newID.String()
	imsu.Cache.GetRWMutex().Lock()
	defer imsu.Cache.GetRWMutex().Unlock()
	imsu.Cache.GetShortURIUserCache()[entity.ID] = entity

	return entity, nil
}

func (imsu *shortURIUserInMemRepo) Change(ctx context.Context, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, _err.NewAppModelValidationError("short_uri_user", err)
	}
	find, err := imsu.Get(ctx, entity.ID)
	if err != nil {
		return nil, err
	}
	if find == nil {
		return nil, _err.NewAppModelNotFoundError(entity.ID, "short_uri_user", "")
	}

	imsu.Cache.GetRWMutex().Lock()
	defer imsu.Cache.GetRWMutex().Unlock()
	imsu.Cache.GetShortURIUserCache()[entity.ID] = entity

	return entity, nil
}

func (imsu *shortURIUserInMemRepo) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	find, err := imsu.Get(ctx, ID)
	if err != nil {
		return err
	}
	if find == nil {
		return nil
	}

	return imsu.delete(ctx, find)
}

func (imsu *shortURIUserInMemRepo) DeleteByUnique(ctx context.Context, userID string, ID string) error {
	if userID == "" || ID == "" {
		return nil
	}

	find, err := imsu.GetByUnique(ctx, userID, ID)
	if err != nil {
		return err
	}
	if find == nil {
		return nil
	}

	return imsu.delete(ctx, find)
}

func (imsu *shortURIUserInMemRepo) DeleteAllByUnique(ctx context.Context, userID string, shortURIIds []string) error {
	// ToDo: implement goroutine
	// ..

	if userID == "" || len(shortURIIds) == 0 {
		return nil
	}

	for _, shortURIID := range shortURIIds {
		err := imsu.DeleteByUnique(ctx, userID, shortURIID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *shortURIUserInMemRepo) DeleteAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	toDelete, err := imsu.ListAllByUser(ctx, userID)
	if err != nil {
		return err
	}
	for _, entity := range toDelete {
		err := imsu.delete(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *shortURIUserInMemRepo) DeleteAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	toDelete, err := imsu.ListAllByShortURI(ctx, shortURIID)
	if err != nil {
		return err
	}

	for _, entity := range toDelete {
		err := imsu.Delete(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *shortURIUserInMemRepo) Remove(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	imsu.Cache.GetRWMutex().Lock()
	defer imsu.Cache.GetRWMutex().Unlock()
	delete(imsu.Cache.GetShortURIUserCache(), ID)

	return nil
}

func (imsu *shortURIUserInMemRepo) RemoveByUnique(ctx context.Context, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	find, err := imsu.GetByUnique(ctx, userID, shortURIID)
	if err != nil {
		return err
	}
	if find == nil {
		return nil
	}

	return imsu.Remove(ctx, find.ID)
}

func (imsu *shortURIUserInMemRepo) RemoveAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	toRemove, err := imsu.ListAllByUser(ctx, userID)
	if err != nil {
		return err
	}

	for _, entity := range toRemove {
		err := imsu.Remove(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *shortURIUserInMemRepo) RemoveAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	toRemove, err := imsu.ListAllByShortURI(ctx, shortURIID)
	if err != nil {
		return err
	}

	for _, entity := range toRemove {
		err := imsu.Remove(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *shortURIUserInMemRepo) CreateStmt(ctx context.Context, stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	return imsu.Create(ctx, entity)
}

func (imsu *shortURIUserInMemRepo) ChangeStmt(ctx context.Context, stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	return imsu.Change(ctx, entity)
}

func (imsu *shortURIUserInMemRepo) DeleteStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	return imsu.Delete(ctx, ID)
}

func (imsu *shortURIUserInMemRepo) DeleteByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	return imsu.DeleteByUnique(ctx, userID, shortURIID)
}

func (imsu *shortURIUserInMemRepo) DeleteAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIds []string) error {
	return imsu.DeleteAllByUnique(ctx, userID, shortURIIds)
}

func (imsu *shortURIUserInMemRepo) DeleteAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	return imsu.DeleteAllByUser(ctx, userID)
}

func (imsu *shortURIUserInMemRepo) DeleteAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	return imsu.DeleteAllByShortURI(ctx, shortURIID)
}

func (imsu *shortURIUserInMemRepo) RemoveStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	return imsu.Remove(ctx, ID)
}

func (imsu *shortURIUserInMemRepo) RemoveByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	return imsu.RemoveByUnique(ctx, userID, shortURIID)
}

func (imsu *shortURIUserInMemRepo) RemoveAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	return imsu.RemoveAllByUser(ctx, userID)
}

func (imsu *shortURIUserInMemRepo) RemoveAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	return imsu.RemoveAllByShortURI(ctx, shortURIID)
}

func (imsu *shortURIUserInMemRepo) delete(ctx context.Context, entity *_model.ShortURIUser) error {
	if entity == nil {
		return nil
	}

	entity.Deleted = true

	_, err := imsu.Change(ctx, entity)

	return err
}
