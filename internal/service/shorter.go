package service

import (
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
	// ..
}

// ================
