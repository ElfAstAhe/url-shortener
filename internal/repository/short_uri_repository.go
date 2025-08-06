package repository

import "github.com/ElfAstAhe/url-shortener/internal/model"

type ShortUriRepository interface {
	GetById(id string) (*model.ShortUri, error)
	GetByKey(key string) (*model.ShortUri, error)
	Create(shortUri *model.ShortUri) (*model.ShortUri, error)
}
