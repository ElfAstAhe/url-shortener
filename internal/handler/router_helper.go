package handler

import (
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_repo "github.com/ElfAstAhe/url-shortener/internal/repository"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
)

func (cr *chiRouter) createShortenService() (_srv.ShorterService, error) {
	db, err := _db.NewDB(cr.config.DBKind, cr.config.DBDsn)
	if err != nil {
		return nil, err
	}
	repository, err := _repo.NewShortURIRepository(db)
	if err != nil {
		return nil, err
	}

	return _srv.NewShorterService(repository)
}

func (cr *chiRouter) createDBConnCheckService() (_repo.DBConnCheckRepository, error) {
	db, err := _db.NewPGIter10Gap(cr.config.DBDsn)
	if err != nil {
		return nil, err
	}

	return _repo.NewDBConnCheckRepository(db)
}
