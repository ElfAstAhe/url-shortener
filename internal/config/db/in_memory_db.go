package db

import "github.com/ElfAstAhe/url-shortener/internal/model"

type InMemoryDB struct {
	ShortURI map[string]*model.ShortURI
}

var InMemoryDBInstance *InMemoryDB

func newInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		ShortURI: make(map[string]*model.ShortURI),
	}
}

func init() {
	InMemoryDBInstance = newInMemoryDB()
}
