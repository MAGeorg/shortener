package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger

type LogStruct struct {
	URI      string
	Method   string
	Status   int
	Duration time.Duration
	Size     int
}

func NewLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer func() {
		// под вопросом куда ставить
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()
	Sugar = *logger.Sugar()

	return nil
}

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
