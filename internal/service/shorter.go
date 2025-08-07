package service

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_repos "github.com/ElfAstAhe/url-shortener/internal/repository"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
)

type Shorter struct {
	Repository _repos.ShortUriRepository
}

func NewShorterService() *Shorter {
	return &Shorter{Repository: _repos.NewShortUriRepository(&_cfg.GlobalConfig.Db)}
}

func (s Shorter) GetUrl(key string) (string, error) {
	model, err := s.Repository.GetByKey(key)
	if err != nil {
		return "", err
	}
	if model == nil {
		return "", nil
	}

	return model.OriginalUrl.String(), nil
}

func (s Shorter) Store(url string) (string, error) {
	key := _utl.EncodeUriStr(url)
	model, err := _model.NewShortUri(url, key)
	if err != nil {
		return "", err
	}

	model, err = s.Repository.Create(model)
	if err != nil {
		return "", err
	}

	return model.Key, nil
}
