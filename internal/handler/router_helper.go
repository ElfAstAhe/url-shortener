package handler

import (
	"net/http"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_repo "github.com/ElfAstAhe/url-shortener/internal/repository"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
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

func (cr *chiRouter) processUnauthorizedIter14(rw http.ResponseWriter, message string) error {
	tokenString, err := _auth.NewJWTStringFromUserInfo(_auth.BuildRandomUserInfo())
	if err != nil {
		return err
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     _auth.CookieName,
		Value:    tokenString,
		SameSite: http.SameSiteStrictMode,
	})

	rw.WriteHeader(http.StatusUnauthorized)
	if message != "" {
		if _, err := rw.Write([]byte(message)); err != nil {
			return err
		}
	}

	return nil
}
