package handler

import (
	"io"
	"net/http"

	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
)

func (cr *chiRouter) pingGetHandler(w http.ResponseWriter, r *http.Request) {
	service, err := cr.createDBConnCheckService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer func() {
		if closer, ok := service.(io.Closer); ok {
			_utl.CloseOnly(closer)
		}
	}()

	err = service.CheckDBConn()
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
