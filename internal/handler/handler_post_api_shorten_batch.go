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

func (cr *chiRouter) shortenBatchPostHandler(rw http.ResponseWriter, r *http.Request) {
	userInfo, err := _auth.UserInfoFromRequestJWT(r)
	if err != nil {
		// Attention!!! For iteration 14 ONLY, remove in future!
		message := fmt.Sprintf("userInfoFromRequestJWT error: [%v]", err)
		cr.log.Error(message)
		if userInfo, err = cr.iter14SetAuthCookie(rw); err != nil {
			message := fmt.Sprintf("process unauthorized error: [%v]", err)
			cr.log.Error(message)
		}
	}

	dec := json.NewDecoder(r.Body)
	var income = make([]*_dto.ShortenBatchCreateItem, 0)
	if err := dec.Decode(&income); err != nil {
		message := fmt.Sprintf("Error deserializing request JSON body: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	serviceBatch, err := _mapper.ShortenBatchFromDto(income)
	if err != nil {
		message := fmt.Sprintf("Error map income data into internal structs: [%s]", err)
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

	ctx := context.WithValue(r.Context(), _auth.ContextUserInfo, userInfo)
	serviceBatchResult, err := service.BatchStore(ctx, serviceBatch)
	if err != nil {
		message := fmt.Sprintf("Error processing batch data: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	respBatch, err := _mapper.ShortenBatchResponseFromKeys(serviceBatchResult)
	if err != nil {
		message := fmt.Sprintf("Error converting batch data: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(respBatch); err != nil {
		message := fmt.Sprintf("Error encoding response as JSON: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	cr.log.Debug("Done batch")
}
