package logger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ElfAstAhe/url-shortener/internal/handler/middleware"
	"github.com/ElfAstAhe/url-shortener/internal/logger"
)

func CustomInfoHTTPLogger(nextHandler http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := middleware.NewCommonResponseWriter(rw)

		nextHandler.ServeHTTP(crw, r)

		duration := time.Since(start)

		logger.Log.Sugar().Infof("uri [%s] method [%s] duration [%v]ms status [%v] size [%v]",
			r.RequestURI, r.Method, strconv.FormatInt(duration.Milliseconds(), 10), crw.Info.StatusCode, crw.Info.Size)
	}

	return http.HandlerFunc(fn)
}
