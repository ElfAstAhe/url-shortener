package helper

import (
	"fmt"
	"io"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_repo "github.com/ElfAstAhe/url-shortener/internal/repository"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
)

func CreateService() _srv.ShorterService {
	return _srv.NewShorterService(_repo.NewShortURIRepository(_cfg.AppConfig.DBKind, _cfg.AppConfig.DB))
}

func CloseOnly(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Printf("Error closing reader: %s\r\n", err)
	}
}
