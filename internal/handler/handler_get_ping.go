package handler

import (
	"context"
	"net/http"
	"time"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
)

func (ar *AppRouter) pingGetHandler(w http.ResponseWriter, r *http.Request) {
	db, err := _db.NewDB(_cfg.DBKindPostgres, _cfg.AppConfig.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = db.GetDB().PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
