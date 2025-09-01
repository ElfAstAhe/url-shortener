package handler

import (
	"fmt"
	"io"
	"net/http"

	_helper "github.com/ElfAstAhe/url-shortener/internal/handler/helper"
	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
)

func rootPOSTHandler(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var err error
	// read income data
	data, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	// store data
	var key string
	key, err = _helper.CreateService().Store(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// prepare outcome data
	newURI, err := _mapper.ResponseFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// outcome data
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(newURI.Result))
	if err != nil {
		fmt.Printf("error writing response [%s]", err.Error())

		return
	}
}
