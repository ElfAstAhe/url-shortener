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

func (cr *chiRouter) userUrlsDeleteHandler(rw http.ResponseWriter, r *http.Request) {
	userInfo, err := _auth.UserInfoFromRequestJWT(r)
	if err != nil {
		message := fmt.Sprintf("User info from JWT is invalid: [%v]", err)
		cr.log.Error(message)

		http.Error(rw, message, http.StatusUnauthorized)

		return
	}

	var dto = make(_dto.ShortenBatchDeleteRequest, 0)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&dto); err != nil {
		message := fmt.Sprintf("error deserialize JSON data: [%v]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	ctx := context.WithValue(context.Background(), _auth.ContextUserInfo, userInfo)

	go cr.batchDeleteAsync(ctx, dto)

	rw.WriteHeader(http.StatusAccepted)

	cr.log.Debug("Done")
}

func (cr *chiRouter) batchDeleteAsync(ctx context.Context, request _dto.ShortenBatchDeleteRequest) {
	// create service instance
	service, err := cr.createShortenService()
	if err != nil {
		cr.log.Error(fmt.Sprintf("Create shorten service failed: [%v]", err))

		return
	}

	// map into service format
	source, err := _mapper.UserBatchDeletesFromDto(request)
	if err != nil {
		cr.log.Error(fmt.Sprintf("Error map batch deletes failed: [%v]", err))

		return
	}

	// do the job
	if err := service.BatchDelete(ctx, source); err != nil {
		cr.log.Error(fmt.Sprintf("Batch deletes failed: [%v]", err))
	}
}
