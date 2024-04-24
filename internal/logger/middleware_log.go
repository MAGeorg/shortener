package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// кастомный тип ResponseData, содержащий статус запроса и размер ответа.
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

// структура логгера.
type Logger struct {
	logger *zap.SugaredLogger
}

// создание нового экземпляра логгера Middleware.
func NewLogMiddleware(lg *zap.SugaredLogger) *Logger {
	return &Logger{
		logger: lg,
	}
}

// метод, реализующий обертку логгера над основной логикой endpoint.
func (lg *Logger) LogMiddleware(h http.HandlerFunc) http.Handler {
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

		lg.logger.Infoln(
			"URI", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size payload", responseData.size,
		)
	}

	return http.HandlerFunc(logFunc)
}
