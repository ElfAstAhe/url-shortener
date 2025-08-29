package handler

import (
	"time"

	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func BuildRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//    router.Use(middleware.Logger)
	router.Use(_log.CustomHInfoHTTPLogger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/{key}", rootGETHandler)
	router.Post("/", rootPOSTHandler)

	return router
}
