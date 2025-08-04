package db

import "github.com/ElfAstAhe/url-shortener/internal/model"

type InMemoryDb struct {
	ShortUri map[string]*model.ShortUri
}

var InMemoryDbInstance *InMemoryDb

func newInMemoryDb() *InMemoryDb {
	return &InMemoryDb{
		ShortUri: make(map[string]*model.ShortUri),
	}
}

func init() {
	InMemoryDbInstance = newInMemoryDb()
}
