package handler

import (
	"net/http"
	"time"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_compress "github.com/ElfAstAhe/url-shortener/internal/handler/middleware/compress"
	_log "github.com/ElfAstAhe/url-shortener/internal/handler/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type chiRouter struct {
	Router *chi.Mux
	config *_cfg.Config
	log    *zap.SugaredLogger
}

func NewRouter(config *_cfg.Config, logger *zap.SugaredLogger) AppRouter {
	appRouter := &chiRouter{
		Router: chi.NewRouter(),
		config: config,
		log:    logger,
	}

	appRouter.buildRoutes(appRouter.Router)

	return appRouter
}

// AppRouter

func (cr *chiRouter) GetRouter() http.Handler {
	return cr.Router
}

// =========

func (cr *chiRouter) buildRoutes(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(_compress.CustomCompress(_compress.DefaultCompressionLevel, _compress.ContentTypeApplicationJSON, _compress.ContentTypeTextHTML))
	router.Use(_compress.CustomDecompress)
	router.Use(_log.CustomInfoHTTPLogger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// root routing
	router.Route("/", func(r chi.Router) {
		r.Get("/{key}", cr.rootGETHandler) // GET   /{key}
		r.Post("/", cr.rootPOSTHandler)    // POST  /
	})

	// ping sub router
	router.Route("/ping", func(r chi.Router) {
		r.Get("/", cr.pingGetHandler) // GET /ping
	})

	// api sub router
	router.Route("/api", func(r chi.Router) {
		// shorten resource sub router
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", cr.shortenPostHandler) // POST  /api/shorten
		})
	})
}
