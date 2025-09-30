package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
)

func (cr *chiRouter) rootPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	userInfo, err := _auth.UserInfoFromRequestJWT(r)
	if err != nil {
		// Attention!!! For iteration 14 ONLY, remove in future!
		message := fmt.Sprintf("userInfoFromRequestJWT error: [%v]", err)
		cr.log.Warn(message)
		if userInfo, err = cr.iter14SetAuthCookie(rw); err != nil {
			message := fmt.Sprintf("iter14 set user info cookie error: [%v]", err)
			cr.log.Error(message)

			http.Error(rw, message, http.StatusInternalServerError)

			return
		}
	}

	// read income data
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	service, err := cr.createShortenService()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	ctx := context.WithValue(r.Context(), _auth.ContextUserInfo, userInfo)
	// store data
	key, conflictErr := service.Store(ctx, string(data))
	if conflictErr != nil && key == "" {
		http.Error(rw, conflictErr.Error(), http.StatusInternalServerError)

		return
	}

	// prepare outcome data
	newURI, err := _mapper.ShortenCreateResponseFromKey(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// outcome data
	rw.Header().Set("Content-Type", "text/plain")
	if conflictErr != nil {
		rw.WriteHeader(http.StatusConflict)
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
	_, err = rw.Write([]byte(newURI.Result))
	if err != nil {
		fmt.Printf("error writing response [%s]", err.Error())

		return
	}
}
