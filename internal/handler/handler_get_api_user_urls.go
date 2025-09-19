package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
)

func (cr *chiRouter) userUrlsHandler(rw http.ResponseWriter, r *http.Request) {
	userInfo, err := _auth.UserInfoFromRequestJWT(r)
	if err != nil {
		// Attention!!! For iteration 14 ONLY, remove in future!
		message := fmt.Sprintf("userInfoFromRequestJWT error: [%v]", err)
		cr.log.Error(message)
		if err := cr.iter14ProcessNoContent(rw, message); err != nil {
			message := fmt.Sprintf("process unauthorized error: [%v]", err)
			cr.log.Error(message)
		}

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

	modelData, err := service.GetAllUserShorts(ctx, userInfo.UserID)
	if err != nil {
		message := fmt.Sprintf("Error getting all user shortens: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	data, err := _mapper.UserShortensFromModel(modelData)
	if err != nil {
		message := fmt.Sprintf("Error map user shortens: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if len(data) == 0 {
		rw.WriteHeader(http.StatusNoContent)
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(data); err != nil {
		message := fmt.Sprintf("Error encoding response as JSON: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	cr.log.Debug("Done")
}
