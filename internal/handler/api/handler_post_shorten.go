package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	_dto "github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	_helper "github.com/ElfAstAhe/url-shortener/internal/handler/helper"
	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
)

func ShortenPostHandler(rw http.ResponseWriter, r *http.Request) {
	log := _log.Log.Sugar()
	dec := json.NewDecoder(r.Body)
	var request _dto.ShortenCreateRequest

	if err := dec.Decode(&request); err != nil {
		message := fmt.Sprintf("Error deserializing request JSON body: [%s]", err)
		log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	key, err := _helper.CreateService().Store(request.URL)
	if err != nil {
		message := fmt.Sprintf("Error storing URL [%s]", err)
		log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	resp, _ := _mapper.ResponseFromKey(key)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		message := fmt.Sprintf("Error encoding response as JSON: [%s]", err)
		log.Error(message)

		http.Error(rw, message, http.StatusInternalServerError)
	}

	log.Debug("Done creating shorten URL")
}
