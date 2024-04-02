package logger

import (
	"net/http"
	"time"
)

func MiddlewareLog(h http.HandlerFunc) http.Handler {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &ResponseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Sugar.Infoln(
			"URI", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size payload", responseData.size,
		)

	}

	return http.HandlerFunc(logFunc)
}
