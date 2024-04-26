// пакет содержащий логику логера.
package logger

import (
	"time"

	"go.uber.org/zap"
)

// структура лога для middleware лога.
type LogStruct struct {
	URI      string
	Method   string
	Status   int
	Duration time.Duration
	Size     int
}

// получение нового экземпляра логгера.
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
