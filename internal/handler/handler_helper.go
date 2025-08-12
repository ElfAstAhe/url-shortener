package handler

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_repo "github.com/ElfAstAhe/url-shortener/internal/repository"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
)

func createService() _srv.ShorterService {
	return _srv.NewShorterService(_repo.NewShortURIRepository(_cfg.AppConfig.DBKind, _cfg.AppConfig.DB))
}
