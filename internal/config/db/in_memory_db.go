package db

import (
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type InMemoryDB struct {
	ShortURI map[string]*_model.ShortURI
}

var InMemoryDBInstance *InMemoryDB

func newInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		ShortURI: make(map[string]*_model.ShortURI),
	}
}

func init() {
	InMemoryDBInstance = newInMemoryDB()
}
