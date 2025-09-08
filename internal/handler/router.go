package handler

import (
	"time"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_compress "github.com/ElfAstAhe/url-shortener/internal/handler/middleware/compress"
	_log "github.com/ElfAstAhe/url-shortener/internal/handler/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type AppRouter struct {
	Router *chi.Mux
	config *_cfg.Config
	log    *zap.SugaredLogger
}

func NewRouter(config *_cfg.Config, logger *zap.SugaredLogger) *AppRouter {
	appRouter := &AppRouter{
		Router: chi.NewRouter(),
		config: config,
		log:    logger,
	}

	appRouter.Router.Use(middleware.RequestID)
	appRouter.Router.Use(middleware.RealIP)
	appRouter.Router.Use(_compress.CustomCompress(_compress.DefaultCompressionLevel, _compress.ContentTypeApplicationJSON, _compress.ContentTypeTextHTML))
	appRouter.Router.Use(_compress.CustomDecompress)
	appRouter.Router.Use(_log.CustomInfoHTTPLogger)
	appRouter.Router.Use(middleware.Recoverer)
	appRouter.Router.Use(middleware.Timeout(60 * time.Second))

	// root routing
	appRouter.Router.Route("/", func(r chi.Router) {
		r.Get("/{key}", appRouter.rootGETHandler) // GET   /{key}
		r.Post("/", appRouter.rootPOSTHandler)    // POST  /
	})

	// ping sub router
	appRouter.Router.Route("/ping", func(r chi.Router) {
		r.Get("/", appRouter.pingGetHandler) // GET /ping
	})

	// api sub router
	appRouter.Router.Route("/api", func(r chi.Router) {
		// shorten resource sub router
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", appRouter.shortenPostHandler) // POST  /api/shorten
		})
	})

	return appRouter
}
