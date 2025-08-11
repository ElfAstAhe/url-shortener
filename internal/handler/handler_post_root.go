package handler

import (
	"fmt"
	"io"
	"net/http"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
)

func rootPOSTHandler(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var err error
	// read income data
	data, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// store data
	var key string
	key, err = _srv.NewShorterService().Store(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// prepare outcome data
	newURI := _utl.BuildNewURI(_cfg.GlobalConfig.BaseURL, key)

	// outcome data
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(newURI))
	if err != nil {
		fmt.Printf("error writing response [%s]", err.Error())

		return
	}
}
