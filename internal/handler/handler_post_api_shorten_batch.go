package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
)

func (cr *chiRouter) shortenBatchPostHandler(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var data _srv.CorrelationUrls = make(_srv.CorrelationUrls)

	if err := dec.Decode(&data); err != nil {
		message := fmt.Sprintf("Error deserializing request JSON body: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	service, err := cr.createShortenService()
	if err != nil {
		message := fmt.Sprintf("Error creating shorten service: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	result, err := service.BatchStore(data)
	if err != nil {
		message := fmt.Sprintf("Error processing batch data: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	resp, err := _mapper.ShortenBatchResponseFromKeys(result)
	if err != nil {
		message := fmt.Sprintf("Error converting batch data: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		message := fmt.Sprintf("Error encoding response as JSON: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	cr.log.Debug("Done batch")
}
