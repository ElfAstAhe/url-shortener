package compress

import "net/http"

func CustomHTTPCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

	})
}
