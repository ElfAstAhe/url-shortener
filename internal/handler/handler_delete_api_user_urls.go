package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	"github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	"github.com/ElfAstAhe/url-shortener/internal/service/auth"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
)

func (cr *chiRouter) userUrlsDeleteHandler(rw http.ResponseWriter, r *http.Request) {
	userInfo, err := auth.UserInfoFromRequestJWT(r)
	if err != nil {
		message := fmt.Sprintf("User info from JWT is invalid: [%v]", err)
		cr.log.Error(message)

		http.Error(rw, message, http.StatusUnauthorized)

		return
	}

	var incomeData = make(dto.ShortenBatchDeleteRequest, 0)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&incomeData); err != nil {
		message := fmt.Sprintf("error deserialize JSON data: [%v]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	ctx := context.WithValue(context.Background(), auth.ContextUserInfo, userInfo)

	go cr.batchDeleteAsync(ctx, incomeData)

	rw.WriteHeader(http.StatusAccepted)

	cr.log.Debug("Done")
}

func (cr *chiRouter) batchDeleteAsync(ctx context.Context, request dto.ShortenBatchDeleteRequest) {
	// create service instance
	service, err := cr.createShortenService()
	if err != nil {
		cr.log.Error(fmt.Sprintf("Create shorten service failed: [%v]", err))

		return
	}
	defer utils.CloseOnly(service)

	// map into service format
	source, err := mapper.UserBatchDeletesFromDto(request)
	if err != nil {
		cr.log.Error(fmt.Sprintf("Error map batch deletes failed: [%v]", err))

		return
	}

	// do the job
	if err := service.BatchDelete(ctx, source); err != nil {
		cr.log.Error(fmt.Sprintf("Batch deletes failed: [%v]", err))
	}
}
