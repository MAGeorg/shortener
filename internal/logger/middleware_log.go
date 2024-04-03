package logger

import (
	"net/http"
	"time"
)

type ResponseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *ResponseData
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size = size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

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
