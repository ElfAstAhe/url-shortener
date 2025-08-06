package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ElfAstAhe/url-shortener/internal/repository"
	"github.com/ElfAstAhe/url-shortener/internal/service"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
)

const RootHandlePath string = "/"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rootPOSTHandler(w, r)

		return
	} else if r.Method == http.MethodGet {
		rootGETHandler(w, r)

		return
	}

	http.Error(w, "Only GET and POST method allowed!", http.StatusBadRequest)
}

func rootGETHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, RootHandlePath)
	if !(len(paths) >= 2) {
		http.Error(w, "No key applied: example [http://localhost:8080/{short_key}]", http.StatusBadRequest)

		return
	}

	key := paths[1]
	fullUrl, err := service.NewShorterService(repository.NewShortUriInMemRepo()).GetUrl(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, fullUrl, http.StatusTemporaryRedirect)
}

// small helper for local tests (remove before PR)
func _(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("Method [%s]\r\n", r.Method)
	body += fmt.Sprintf("HEADERS ========================\r\n")
	for k, v := range r.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += fmt.Sprintf("PATH ===========================\r\n")
	body += fmt.Sprintf("Path [%s]\r\n", r.URL.Path)
	body += fmt.Sprintf("Path trimmed [%s]\r\n", strings.TrimPrefix(r.URL.Path, RootHandlePath))
	paths := strings.Split(r.URL.Path, "/")
	body += fmt.Sprintf("Paths array [%v]\r\n", paths)
	body += fmt.Sprintf("QUERY PARAMS ===================\r\n")
	for k, v := range r.URL.Query() {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += fmt.Sprintf("FORM ===========================\r\n")
	for k, v := range r.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += fmt.Sprintf("PATH PARAMS ====================\r\n")
	key := r.PathValue("key")
	if key == "" {
		body += fmt.Sprintf("No {key} param\r\n")
	} else {
		body += fmt.Sprintf("Key [%s]", key)
	}

	w.Write([]byte(body))
}

func rootPOSTHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	hashBytes := utils.EncodeUri(data)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	newUri := utils.BuildNewUri(r.URL, string(hashBytes[:]))

	_, err = w.Write([]byte(newUri))
	if err != nil {
		fmt.Printf("error writing response [%s]", err.Error())

		return
	}
}
