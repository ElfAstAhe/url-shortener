package handler

import (
	"fmt"
	"github.com/ElfAstAhe/url-shortener/internal/utils"
	"io/ioutil"
	"net/http"
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
	http.Error(w, "Method implementation in process", http.StatusMethodNotAllowed)
}

func rootPOSTHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
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
