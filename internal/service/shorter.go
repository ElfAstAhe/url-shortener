package service

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_repo "github.com/ElfAstAhe/url-shortener/internal/repository"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
)

type Shorter struct {
	Repository _repo.ShortURIRepository
}

func NewShorterService(repo _repo.ShortURIRepository) (ShorterService, error) {
	return &Shorter{
		Repository: repo,
	}, nil
}

// ShorterService

func (s *Shorter) GetURL(key string) (string, error) {
	model, err := s.Repository.GetByKey(key)
	if err != nil {
		return "", err
	}
	if model == nil {
		return "", nil
	}

	return model.OriginalURL.URL.String(), nil
}

func (s *Shorter) Store(url string) (string, error) {
	key := _utl.EncodeURIStr(url)
	model, err := _model.NewShortURI(url, key)
	if err != nil {
		return "", err
	}

	model, err = s.Repository.Create(model)
	if err != nil {
		return "", err
	}

	return model.Key, nil
}

func (s *Shorter) BatchStore(source CorrelationUrls) (CorrelationShorts, error) {
	if len(source) == 0 {
		return CorrelationShorts{}, nil
	}

	batch, err := toBatchSource(source)
	if err != nil {
		return nil, err
	}

	batchRes, err := s.Repository.BatchCreate(batch)
	if err != nil {
		return nil, err
	}

	res, err := toBatchResult(batchRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ================

func toBatchSource(source CorrelationUrls) (map[string]*_model.ShortURI, error) {
	batch := make(map[string]*_model.ShortURI)
	for correlation, origURL := range source {
		item, err := _model.NewShortURI(origURL, _utl.EncodeURIStr(origURL))
		if err != nil {
			return nil, err
		}
		batch[correlation] = item
	}

	return batch, nil
}

func toBatchResult(source map[string]*_model.ShortURI) (CorrelationShorts, error) {
	batch := make(CorrelationShorts)
	for correlation, shortURL := range source {
		batch[correlation] = _utl.BuildNewURI(_cfg.AppConfig.BaseURL, shortURL.Key)
	}

	return batch, nil
}
