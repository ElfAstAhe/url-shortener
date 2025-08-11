package handler

import (
	"github.com/go-chi/chi/v5"
)

func BuildRouter() chi.Router {
	router := chi.NewRouter()

	router.Get("/{key}", rootGETHandler)
	router.Post("/", rootPOSTHandler)

	return router
}
