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

func (cr *chiRouter) iter14ProcessUnauthorized(rw http.ResponseWriter, message string) error {
	if _, err := cr.iter14SetAuthCookie(rw); err != nil {
		return err
	}

	rw.WriteHeader(http.StatusUnauthorized)
	if message != "" {
		if _, err := rw.Write([]byte(message)); err != nil {
			return err
		}
	}

	return nil
}

func (cr *chiRouter) iter14ProcessNoContent(rw http.ResponseWriter, message string) error {
	if _, err := cr.iter14SetAuthCookie(rw); err != nil {
		return err
	}

	rw.WriteHeader(http.StatusNoContent)
	if message != "" {
		if _, err := rw.Write([]byte(message)); err != nil {
			return err
		}
	}

	return nil
}

func (cr *chiRouter) iter14SetAuthCookie(rw http.ResponseWriter) (*_auth.UserInfo, error) {
	userInfo := _auth.BuildRandomUserInfo()
	tokenString, err := _auth.NewJWTStringFromUserInfo(userInfo)
	if err != nil {
		return nil, err
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     _auth.CookieName,
		Value:    tokenString,
		SameSite: http.SameSiteStrictMode,
	})

	return userInfo, nil
}
