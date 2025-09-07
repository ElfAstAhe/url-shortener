package handler

import (
	"net/http"

	_helper "github.com/ElfAstAhe/url-shortener/internal/handler/helper"
	"github.com/go-chi/chi/v5"
)

func rootGETHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "No key applied: example [http://localhost:8080/{key}]", http.StatusBadRequest)

		return
	}
	fullURL, err := _helper.CreateService().GetURL(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if fullURL == "" {
		http.Error(w, "No shorter url found!", http.StatusNotFound)

		return
	}

	http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
}
