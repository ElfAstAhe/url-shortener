package logger

import (
	"net/http"
	"strconv"
	"time"
)

func CustomHInfoHTTPLogger(nextHandler http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := NewResponseLoggerWriter(rw)

		nextHandler.ServeHTTP(lrw, r)

		duration := time.Since(start)

		Log.Sugar().Infof("uri [%s] method [%s] duration [%v] status [%v] size [%v]",
			r.RequestURI, r.Method, strconv.FormatInt(duration.Milliseconds(), 10),
			lrw.info.StatusCode, lrw.info.Size)
	}

	return http.HandlerFunc(fn)
}
