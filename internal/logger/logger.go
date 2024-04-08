package logger

import (
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

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
