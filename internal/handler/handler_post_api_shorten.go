package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	_dto "github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
)

func (cr *chiRouter) shortenPostHandler(rw http.ResponseWriter, r *http.Request) {
	userInfo, err := _auth.UserInfoFromRequestJWT(r)
	if err != nil {
		// Attention!!! For iteration 14 ONLY, remove in future!
		message := fmt.Sprintf("userInfoFromRequestJWT error: [%v]", err)
		cr.log.Error(message)
		if err := cr.processUnauthorizedIter14(rw, message); err != nil {
			message := fmt.Sprintf("process unauthorized error: [%v]", err)
			cr.log.Error(message)
		}

		return
	}

	dec := json.NewDecoder(r.Body)
	var request _dto.ShortenCreateRequest
	if err := dec.Decode(&request); err != nil {
		message := fmt.Sprintf("Error deserializing request JSON body: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	if request.URL == "" {
		message := fmt.Sprintf("Empty URL: [%s]", request)
		cr.log.Warn(message)
		http.Error(rw, message, http.StatusBadRequest)

		return
	}

	service, err := cr.createShortenService()
	if err != nil {
		message := fmt.Sprintf("Error creating shorten service: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	ctx := context.WithValue(r.Context(), _auth.ContextUserInfo, userInfo)
	key, conflictErr := service.Store(ctx, request.URL)
	if conflictErr != nil && key == "" {
		message := fmt.Sprintf("Error storing URL [%s]", conflictErr)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	resp, _ := _mapper.ShortenCreateResponseFromKey(key)

	rw.Header().Set("Content-Type", "application/json")
	if conflictErr != nil {
		rw.WriteHeader(http.StatusConflict)
	} else {
		rw.WriteHeader(http.StatusCreated)
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		message := fmt.Sprintf("Error encoding response as JSON: [%s]", err)
		cr.log.Error(message)

		http.Error(rw, message, http.StatusInternalServerError)
	}

	cr.log.Debug("Done creating shorten URL")
}
