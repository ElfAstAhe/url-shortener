package utils

import (
	"net/http"
	"net/url"
)

// BuildNewUri - build URI for new resource
func BuildNewUri(request *http.Request, key string) string {
	// current using only http scheme for dev iteration 1
	// new URL
	newURL := &url.URL{
		Scheme: "http",
		Host:   request.Host,
	}

	// add path
	newURL = newURL.JoinPath("/", key)

	return newURL.String()
}
