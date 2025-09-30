package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
	_err "github.com/ElfAstAhe/url-shortener/pkg/errors"
	"github.com/go-chi/chi/v5"
)

func (cr *chiRouter) rootGETHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := _auth.UserInfoFromRequestJWT(r)
	if err != nil {
		message := fmt.Sprintf("userInfoFromRequestJWT error: [%v]", err)
		cr.log.Error(message)
		http.Error(w, message, http.StatusUnauthorized)

		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "No key applied: example [http://localhost:8080/{key}]", http.StatusBadRequest)

		return
	}
	service, err := cr.createShortenService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	ctx := context.WithValue(r.Context(), _auth.ContextUserInfo, userInfo)
	fullURL, err := service.GetURL(ctx, key)
	if err != nil {
		if errors.As(err, &_err.AppSoftRemoved) {
			cr.log.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusGone)

			return
		}
		cr.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if fullURL == "" {
		http.Error(w, "No shorter url found!", http.StatusNotFound)

		return
	}

	http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
}
