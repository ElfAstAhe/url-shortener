package service

import (
	"context"
	"errors"
	"io"

	"github.com/ElfAstAhe/url-shortener/internal/config"
	"github.com/ElfAstAhe/url-shortener/internal/model"
	"github.com/ElfAstAhe/url-shortener/internal/repository"
	"github.com/ElfAstAhe/url-shortener/internal/service/auth"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
	_errs "github.com/ElfAstAhe/url-shortener/pkg/errors"
)

type Shorter struct {
	Repository repository.ShortURIRepository
}

func NewShorterService(repo repository.ShortURIRepository) (*Shorter, error) {
	return &Shorter{
		Repository: repo,
	}, nil
}

// Closer

func (s *Shorter) Close() error {
	if closer, ok := s.Repository.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

// ShorterService

func (s *Shorter) GetURL(ctx context.Context, key string) (string, error) {
	noAuthData, errNoAuth := s.getURLNoAuth(ctx, key)
	_, errAuth := s.getURLAuth(ctx, key)
	if errAuth != nil && errors.As(errAuth, &_errs.AppSoftRemoved) {
		return "", errAuth
	}
	if errNoAuth != nil {
		return "", errNoAuth
	}

	return noAuthData, nil
}

func (s *Shorter) getURLNoAuth(ctx context.Context, key string) (string, error) {
	res, err := s.Repository.GetByKey(ctx, key)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.OriginalURL.URL.String(), nil
}

func (s *Shorter) getURLAuth(ctx context.Context, key string) (string, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", _errs.NewAppAuthInfoAbsentError("getURLAuth internal service method", nil)
	}
	res, err := s.Repository.GetByKeyUser(ctx, userInfo.UserID, key)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.OriginalURL.URL.String(), nil
}

func (s *Shorter) Store(ctx context.Context, url string) (string, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return "", err
	}

	key := utils.EncodeURIStr(url)
	res, err := model.NewShortURI(url, key)
	if err != nil {
		return "", err
	}

	res, err = s.Repository.Create(ctx, userInfo.UserID, res)
	if err != nil && res == nil {
		return "", err
	} else if err != nil {
		return res.Key, err
	}

	return res.Key, nil
}

func (s *Shorter) BatchStore(ctx context.Context, source CorrelationUrls) (CorrelationShorts, error) {
	if len(source) == 0 {
		return CorrelationShorts{}, nil
	}

	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	repoBatch, err := s.toBatchSource(source)
	if err != nil {
		return nil, err
	}

	batchRes, err := s.Repository.BatchCreate(ctx, userInfo.UserID, repoBatch)
	if err != nil {
		return nil, err
	}

	res, err := s.toBatchResult(batchRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Shorter) GetAllUserShorts(ctx context.Context, userID string) (UserShorts, error) {
	entities, err := s.Repository.ListAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	models, err := s.toUserShorts(entities)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (s *Shorter) BatchDelete(ctx context.Context, data UserBatchDeletes) error {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return err
	}

	return s.Repository.BatchDeleteByKeys(ctx, userInfo.UserID, data)
}

// ================

func (s *Shorter) toBatchSource(source CorrelationUrls) (map[string]*model.ShortURI, error) {
	batch := make(map[string]*model.ShortURI)
	for correlation, origURL := range source {
		item, err := model.NewShortURI(origURL, utils.EncodeURIStr(origURL))
		if err != nil {
			return nil, err
		}
		batch[correlation] = item
	}

	return batch, nil
}

func (s *Shorter) toBatchResult(source map[string]*model.ShortURI) (CorrelationShorts, error) {
	batch := make(CorrelationShorts)
	for correlation, shortURL := range source {
		batch[correlation] = utils.BuildNewURI(config.AppConfig.BaseURL, shortURL.Key)
	}

	return batch, nil
}

func (s *Shorter) toUserShorts(entities []*model.ShortURI) (UserShorts, error) {
	if len(entities) == 0 {
		return nil, nil
	}
	res := make(UserShorts)
	for _, entity := range entities {
		res[entity.OriginalURL.URL.String()] = utils.BuildNewURI(config.AppConfig.BaseURL, entity.Key)
	}

	return res, nil
}
