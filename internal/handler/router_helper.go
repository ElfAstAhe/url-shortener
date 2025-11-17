package handler

import (
	"net/http"

	"github.com/ElfAstAhe/url-shortener/internal/config/db"
	"github.com/ElfAstAhe/url-shortener/internal/repository"
	"github.com/ElfAstAhe/url-shortener/internal/service"
	"github.com/ElfAstAhe/url-shortener/internal/service/auth"
)

func (cr *chiRouter) createShortenService() (service.ShorterService, error) {
	appDb, err := db.NewDB(cr.config.DBKind, cr.config.DBDsn)
	if err != nil {
		return nil, err
	}
	repo, err := repository.NewShortURIRepository(appDb)
	if err != nil {
		return nil, err
	}

	return service.NewShorterService(repo)
}

func (cr *chiRouter) createDBConnCheckService() (repository.DBConnCheckRepository, error) {
	appDb, err := db.NewPGIter10Gap(cr.config.DBDsn)
	if err != nil {
		return nil, err
	}

	return repository.NewDBConnCheckRepository(appDb)
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

func (cr *chiRouter) iter14SetAuthCookie(rw http.ResponseWriter) (*auth.UserInfo, error) {
	userInfo := auth.BuildRandomUserInfo()
	tokenString, err := auth.NewJWTStringFromUserInfo(userInfo)
	if err != nil {
		return nil, err
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     auth.CookieName,
		Value:    tokenString,
		SameSite: http.SameSiteStrictMode,
	})

	return userInfo, nil
}
