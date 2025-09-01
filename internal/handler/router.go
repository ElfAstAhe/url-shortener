package handler

import (
	"time"

	_api "github.com/ElfAstAhe/url-shortener/internal/handler/api"
	_compress "github.com/ElfAstAhe/url-shortener/internal/handler/middleware/compress"
	_log "github.com/ElfAstAhe/url-shortener/internal/handler/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func BuildRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//    router.Use(middleware.Compress(5))
	router.Use(_compress.CustomCompress(_compress.DefaultCompressionLevel, _compress.ContentTypeApplicationJSON, _compress.ContentTypeTextHTML))
	//    router.Use(middleware.Logger)
	router.Use(_log.CustomInfoHTTPLogger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// root routing
	router.Get("/{key}", rootGETHandler) // GET   /{key}
	router.Post("/", rootPOSTHandler)    // POST  /

	// api sub router
	router.Route("/api", func(r chi.Router) {
		// shorten resource sub router
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", _api.ShortenPostHandler) // POST  /api/shorten
		})
	})

	return router
}
