package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (ar *AppRouter) rootGETHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "No key applied: example [http://localhost:8080/{key}]", http.StatusBadRequest)

		return
	}
	service, err := ar.createShortenService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	fullURL, err := service.GetURL(key)
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
