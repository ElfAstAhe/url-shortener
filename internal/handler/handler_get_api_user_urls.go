package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	_mapper "github.com/ElfAstAhe/url-shortener/internal/handler/mapper"
	_auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"
)

func (cr *chiRouter) userUrlsHandler(rw http.ResponseWriter, r *http.Request) {
	userID, err := _auth.GetUserID(r)
	if err != nil {
		message := fmt.Sprintf("error getting user id: [%v]", err)
		cr.log.Error(message)

		if err := processUnauthorized(rw, message); err != nil {
			message := fmt.Sprintf("error processing unauthorized request: [%v]", err)
			cr.log.Error(message)
			http.Error(rw, message, http.StatusInternalServerError)
		}

		return
	}

	if userID == "" {
		if err := processUnauthorized(rw, ""); err != nil {
			message := fmt.Sprintf("error processing empty user id request: [%v]", err)
			cr.log.Error(message)
			http.Error(rw, message, http.StatusInternalServerError)
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

	modelData, err := service.GetAllUserShorts(userID)
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
	if len(data) == 0 {
		rw.WriteHeader(http.StatusNoContent)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(data); err != nil {
		message := fmt.Sprintf("Error encoding response as JSON: [%s]", err)
		cr.log.Error(message)
		http.Error(rw, message, http.StatusInternalServerError)

		return
	}

	cr.log.Debug("Done")
}

func processUnauthorized(rw http.ResponseWriter, message string) error {
	tokenString, err := _auth.NewTokenString(_auth.TestUser, _auth.TestUserID, _auth.TestRoles...)
	if err != nil {
		return err
	}

	http.SetCookie(rw, &http.Cookie{
		Name:  _auth.Cookie,
		Value: tokenString,
	})

	rw.WriteHeader(http.StatusUnauthorized)
	if message != "" {
		if _, err := rw.Write([]byte(message)); err != nil {
			return err
		}
	}

	return nil
}
