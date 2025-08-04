package service

import (
	"github.com/ElfAstAhe/url-shortener/internal/repository"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
)

type Shorter struct {
	Repository repository.ShortUriRepository
}

func NewShorterService(repository repository.ShortUriRepository) *Shorter {
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
	key := utils.EncodeUriStr(url)

}
