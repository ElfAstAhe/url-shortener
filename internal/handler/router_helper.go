package handler

import (
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_repo "github.com/ElfAstAhe/url-shortener/internal/repository"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
)

func (ar *AppRouter) createShortenService() (_srv.ShorterService, error) {
	db, err := _db.NewDB(ar.config.DBKind, ar.config.DB)
	if err != nil {
		return nil, err
	}
	repository, err := _repo.NewShortURIRepository(db)
	if err != nil {
		return nil, err
	}

	return _srv.NewShorterService(repository)
}
