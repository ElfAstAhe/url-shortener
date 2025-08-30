package handler

import (
	"fmt"
	"io"
	"net/http"

	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_helper "github.com/ElfAstAhe/url-shortener/internal/handler/helper"
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
	key, err = _helper.CreateService().Store(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// prepare outcome data
	newURI := _utl.BuildNewURI(_cfg.AppConfig.BaseURL, key)

	// outcome data
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(newURI))
	if err != nil {
		fmt.Printf("error writing response [%s]", err.Error())

		return
	}
}
