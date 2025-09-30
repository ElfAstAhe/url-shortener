package repository

import (
	"context"
	"errors"
	"sync"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type shortURIInMemRepo struct {
	Cache    _db.InMemoryCache
	userRepo ShortURIUserRepository
	anchor   sync.RWMutex
	log      *zap.SugaredLogger
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
			log:      _log.Log.Sugar(),
		}, nil
	}

	return nil, _err.NewAppInvalidArgument("db param does not implement InMemoryCache")
}

func (ims *shortURIInMemRepo) Get(ctx context.Context, id string) (*_model.ShortURI, error) {
	ims.anchor.RLock()
	defer ims.anchor.RUnlock()
	res := ims.Cache.GetShortURICache()[id]

	return res, nil
}

func (ims *shortURIInMemRepo) GetByKey(ctx context.Context, key string) (*_model.ShortURI, error) {
	if key == "" {
		return nil, nil
	}

	ims.anchor.RLock()
	defer ims.anchor.RUnlock()
	for _, value := range ims.Cache.GetShortURICache() {
		if value.Key == key {
			return value, nil
		}
	}

	return nil, nil
}

func (ims *shortURIInMemRepo) GetByKeyUser(ctx context.Context, userID string, key string) (*_model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	if key == "" {
		return nil, nil
	}

	entity, err := ims.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	userLink, err := ims.userRepo.GetByUnique(ctx, userID, entity.ID)
	if err != nil {
		return nil, err
	}
	if userLink == nil {
		return nil, nil
	}
	if userLink.Deleted {
		return nil, _err.NewAppSoftRemovedError("short_uri", nil)
	}

	return entity, nil
}

func (ims *shortURIInMemRepo) Create(ctx context.Context, userID string, entity *_model.ShortURI) (*_model.ShortURI, error) {
	if err := _model.ValidateShortURI(entity); err != nil {
		return nil, err
	}

	find, err := ims.GetByKey(ctx, entity.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		if err := ims.addUser(ctx, find.ID, userID); err != nil {
			return nil, err
		}

		return find, errors.New("short URI already exists")
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	ims.anchor.Lock()
	defer ims.anchor.Unlock()
	ims.Cache.GetShortURICache()[entity.ID] = entity

	if err := ims.addUser(ctx, entity.ID, userID); err != nil {
		return nil, err
	}

	return entity, nil
}

func (ims *shortURIInMemRepo) BatchCreate(ctx context.Context, userID string, batch map[string]*_model.ShortURI) (map[string]*_model.ShortURI, error) {
	res := make(map[string]*_model.ShortURI)
	if len(batch) == 0 {
		return res, nil
	}

	for correlation, item := range batch {
		data, err := ims.Create(ctx, userID, item)
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

func (ims *shortURIInMemRepo) ListAllByKeys(ctx context.Context, keys []string) ([]*_model.ShortURI, error) {
	res := make([]*_model.ShortURI, 0)
	if len(keys) == 0 {
		return res, nil
	}

	ims.anchor.RLock()
	defer ims.anchor.RUnlock()
	for key, value := range ims.Cache.GetShortURICache() {
		if value.Key == key {
			res = append(res, value)
		}
	}

	return res, nil
}

func (ims *shortURIInMemRepo) Delete(ctx context.Context, ID string, userID string) error {
	if userID == "" {
		return nil
	}
	if ID == "" {
		return nil
	}

	return ims.userRepo.DeleteByUnique(ctx, userID, ID)
}

func (ims *shortURIInMemRepo) BatchDeleteByKeys(ctx context.Context, userID string, keys []string) error {
	if userID == "" || len(keys) == 0 {
		return nil
	}
	ids, err := ims.listIdsByKeys(ctx, keys)
	if err != nil {
		return err
	}

	inCh := ims.iter15Generator(ctx, ids)

	channels := ims.iter15FanOut(ctx, userID, inCh)

	res := ims.iter15FanIn(ctx, channels...)

	var errs errgroup.Group
	for dml := range res {
		errs.Go(func() error {
			ims.log.Infof("BatchDeleteByKeys dml result: [%v]", dml)

			return dml.Err
		})

	}

	return errs.Wait()
}

func (ims *shortURIInMemRepo) iter15Generator(ctx context.Context, ids []string) chan string {
	inCh := make(chan string)

	go func() {
		defer close(inCh)

		for _, id := range ids {
			select {
			case <-ctx.Done():
				return
			case inCh <- id:
			}
		}
	}()

	return inCh
}

func (ims *shortURIInMemRepo) iter15FanOut(ctx context.Context, userID string, inCh chan string) []chan *DMLResult {
	workersCount := 4
	res := make([]chan *DMLResult, workersCount)

	for index := 0; index < workersCount; index++ {
		res[index] = ims.iter15Delete(ctx, userID, inCh)
	}

	return res
}

func (ims *shortURIInMemRepo) iter15Delete(ctx context.Context, userID string, inCh chan string) chan *DMLResult {
	res := make(chan *DMLResult)

	go func() {
		defer close(res)

		for id := range inCh {
			select {
			case <-ctx.Done():
				return
			case res <- NewDMLResult(ims.userRepo.DeleteByUnique(ctx, userID, id), "short_uri_users", id):
			}
		}
	}()

	return res
}

func (ims *shortURIInMemRepo) iter15FanIn(ctx context.Context, resCh ...chan *DMLResult) chan *DMLResult {
	res := make(chan *DMLResult)

	var wg sync.WaitGroup

	for _, ch := range resCh {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case res <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()

		close(res)
	}()

	return res
}

func (ims *shortURIInMemRepo) listIdsByKeys(ctx context.Context, keys []string) ([]string, error) {
	res := make([]string, 0)
	if len(keys) == 0 {
		return res, nil
	}

	for _, key := range keys {
		entity, err := ims.GetByKey(ctx, key)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			continue
		}
		res = append(res, entity.ID)
	}

	return res, nil
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
