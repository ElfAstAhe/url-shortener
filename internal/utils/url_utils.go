package utils

import (
	"net/url"
)

// BuildNewUri - build URI for new resource
func BuildNewUri(reqUrl *url.URL, newPath string) string {
	// new URL
	newUrl := url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	// add path
	newUrl.JoinPath(newPath)

	return newUrl.String()
}
