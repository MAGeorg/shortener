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

func NewLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	Sugar = *logger.Sugar()

	return nil
}
