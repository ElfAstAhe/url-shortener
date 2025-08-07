package utils

import (
	"net/http"
	"net/url"
)

// BuildNewUri - build URI for new resource
func BuildNewUri(request *http.Request, key string) string {
	// current using only http scheme for dev iteration 1
	// new URL
	newUrl := &url.URL{
		Scheme: "http",
		Host:   request.Host,
	}

	// add path
	newUrl = newUrl.JoinPath("/", key)

	return newUrl.String()
}
