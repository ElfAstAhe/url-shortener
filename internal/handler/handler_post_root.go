package handler

import (
	"fmt"
	"io"
	"net/http"

	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
)

func (cr *chiRouter) rootPOSTHandler(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var err error
	// read income data
	data, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	service, err := cr.createShortenService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// store data
	key, conflictErr := service.Store(string(data))
	if conflictErr != nil && key == "" {
		http.Error(w, conflictErr.Error(), http.StatusInternalServerError)

		return
	}

	// prepare outcome data
	newURI, err := _mapper.ShortenCreateResponseFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// outcome data
	w.Header().Set("Content-Type", "text/plain")
	if conflictErr != nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	_, err = w.Write([]byte(newURI.Result))
	if err != nil {
		fmt.Printf("error writing response [%s]", err.Error())

		return
	}
}
