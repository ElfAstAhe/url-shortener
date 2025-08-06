package service

import (
	models "github.com/ElfAstAhe/url-shortener/internal/model"
	repos "github.com/ElfAstAhe/url-shortener/internal/repository"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
)

type Shorter struct {
	Repository repos.ShortUriRepository
}

func NewShorterService(repository repos.ShortUriRepository) *Shorter {
	return &Shorter{Repository: repository}
}

func (s Shorter) GetUrl(key string) (string, error) {
	model, err := s.Repository.GetByKey(key)
	if err != nil {
		return "", err
	}

	return model.OriginalUrl.String(), nil
}

func (s Shorter) Store(url string) (string, error) {
	key := string(utils.EncodeUriStr(url))
	model, err := models.NewShortUri(url, key)
	if err != nil {
		return "", err
	}

	model, err = s.Repository.Create(model)
	if err != nil {
		return "", err
	}

	return model.Key, nil
}
