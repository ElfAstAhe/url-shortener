package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func BuildNewURI(baseURL string, key string) string {
	var _baseURL string = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	var _key string = strings.TrimSpace(key)
	if _baseURL == "" || _key == "" {
		return ""
	}
	_resURL, err := url.Parse(_baseURL)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return _resURL.JoinPath(_key).String()
}
